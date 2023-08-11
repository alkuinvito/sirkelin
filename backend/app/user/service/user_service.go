package service

import (
	"errors"
	"gorm.io/gorm"
	"sirkelin/backend/initializers"
	"sirkelin/backend/models"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	"sirkelin/backend/app/user/repository"

	"firebase.google.com/go/auth"
)

const EXPIRES_IN = time.Hour * 24

type UserService struct {
	repository *repository.UserRepository
	db         *gorm.DB
}

type IUserService interface {
	GetAll() ([]models.User, error)
	GetByID(uid string) (*models.User, error)
	getSessionToken(c *gin.Context, client *auth.Client) (*auth.Token, error)
	initClient(c *gin.Context) (*auth.Client, error)
	revokeToken(c *gin.Context, client *auth.Client) error
	SignIn(c *gin.Context, tokenString string) (*models.User, string, error)
	SignOut(c *gin.Context) error
	UpdateProfile(id string, data models.UpdateProfileSchema) error
	verifyIDToken(c *gin.Context, client *auth.Client, tokenString string) (*auth.Token, error)
	VerifySessionToken(c *gin.Context) (*auth.Token, error)
}

func NewUserService(repository *repository.UserRepository, db *gorm.DB) *UserService {
	return &UserService{
		repository: repository,
		db:         db,
	}
}

func (service *UserService) GetAll() ([]models.User, error) {
	tx := service.db
	users, err := service.repository.Get(tx)
	if err != nil {
		return []models.User{}, err
	}
	return users, nil
}

func (service *UserService) GetByID(uid string) (*models.User, error) {
	tx := service.db
	user, err := service.repository.GetByID(tx, uid)
	if err != nil {
		return &models.User{}, err
	}
	return user, nil
}

func (service *UserService) getSessionToken(c *gin.Context, client *auth.Client) (*auth.Token, error) {
	bearerToken := c.GetHeader("Authorization")
	tokenString := strings.Split(bearerToken, " ")

	if len(tokenString) != 2 {
		cookies, err := c.Cookie("session")
		if err != nil {
			return nil, err
		}
		return client.VerifySessionCookieAndCheckRevoked(c, cookies)
	}
	return client.VerifySessionCookieAndCheckRevoked(c, tokenString[1])
}

func (service *UserService) initClient(c *gin.Context) (*auth.Client, error) {
	return initializers.InitializeAppDefault().Auth(c)
}

func (service *UserService) revokeToken(c *gin.Context, client *auth.Client) error {
	token, err := service.getSessionToken(c, client)
	if err != nil {
		return err
	}

	if err := client.RevokeRefreshTokens(c, token.UID); err != nil {
		return err
	}

	return nil
}

func (service *UserService) SignIn(c *gin.Context, tokenString string) (*models.User, string, error) {
	client, err := service.initClient(c)
	if err != nil {
		return &models.User{}, "", err
	}

	token, err := service.verifyIDToken(c, client, tokenString)
	if err != nil {
		return &models.User{}, "", err
	}

	tx := service.db.Begin()
	defer initializers.CommitOrRollback(tx)
	user, err := service.repository.GetByID(tx, token.Subject)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			err = service.repository.Save(tx,
				&models.User{
					ID:       token.Subject,
					Fullname: token.Claims["name"].(string),
					Picture:  token.Claims["picture"].(string),
					Email:    token.Claims["email"].(string),
				},
			)

			if err != nil {
				return &models.User{}, "", err
			}
		} else {
			return &models.User{}, "", err
		}
	}

	if err != nil {
		return &models.User{}, "", err
	}

	session, err := client.SessionCookie(c, tokenString, EXPIRES_IN)
	return user, session, err
}

func (service *UserService) SignOut(c *gin.Context) error {
	client, err := service.initClient(c)
	if err != nil {
		return err
	}
	return service.revokeToken(c, client)
}

func (service *UserService) UpdateProfile(id string, data models.UpdateProfileSchema) error {
	updated := &models.User{ID: id, Username: data.Username, Fullname: data.Fullname, Picture: data.Picture}
	tx := service.db.Begin()
	defer initializers.CommitOrRollback(tx)
	return service.repository.Update(tx, updated)
}

func (service *UserService) verifyIDToken(c *gin.Context, client *auth.Client, tokenString string) (*auth.Token, error) {
	token, err := client.VerifyIDTokenAndCheckRevoked(c, tokenString)
	if err != nil {
		if err.Error() == "ID token has been revoked" {
			return nil, errors.New("user must reauthenticate")
		} else {
			return nil, errors.New("invalid id token")
		}
	}

	if time.Now().Unix()-token.IssuedAt > 5*60 {
		return nil, errors.New("recent sign-in required")
	}

	return token, err
}

func (service *UserService) VerifySessionToken(c *gin.Context) (*auth.Token, error) {
	client, err := service.initClient(c)
	if err != nil {
		return nil, err
	}
	return service.getSessionToken(c, client)
}
