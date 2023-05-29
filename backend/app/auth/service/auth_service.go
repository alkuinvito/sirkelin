package service

import (
	"errors"
	"gorm.io/gorm"
	"sirkelin/backend/initializers"
	"sirkelin/backend/models"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	"sirkelin/backend/app/auth/repository"

	"firebase.google.com/go/auth"
)

const EXPIRES_IN = time.Hour * 24

type AuthService struct {
	repository *repository.AuthRepository
	db         *gorm.DB
}

type IAuthService interface {
	getSessionToken(c *gin.Context, client *auth.Client) (*auth.Token, error)
	initClient(c *gin.Context) (*auth.Client, error)
	revokeToken(c *gin.Context, client *auth.Client) error
	SignIn(c *gin.Context, tokenString string) (string, error)
	SignOut(c *gin.Context) error
	verifyIDToken(c *gin.Context, client *auth.Client, tokenString string) (*auth.Token, error)
	VerifySessionToken(c *gin.Context) (*auth.Token, error)
}

func NewAuthService(repository *repository.AuthRepository, db *gorm.DB) *AuthService {
	return &AuthService{
		repository: repository,
		db:         db,
	}
}

func (service *AuthService) getSessionToken(c *gin.Context, client *auth.Client) (*auth.Token, error) {
	bearerToken := c.GetHeader("Authorization")
	tokenString := strings.Split(bearerToken, " ")

	if len(tokenString) != 2 {
		cookies, err := c.Cookie("session")
		if err != nil {
			return nil, err
		}
		return client.VerifySessionCookieAndCheckRevoked(c, cookies)
	}
	return client.VerifySessionCookieAndCheckRevoked(c, tokenString[0])
}

func (service *AuthService) initClient(c *gin.Context) (*auth.Client, error) {
	return initializers.InitializeAppDefault().Auth(c)
}

func (service *AuthService) revokeToken(c *gin.Context, client *auth.Client) error {
	token, err := service.getSessionToken(c, client)
	if err != nil {
		return err
	}

	if err := client.RevokeRefreshTokens(c, token.UID); err != nil {
		return err
	}

	return nil
}

func (service *AuthService) SignIn(c *gin.Context, tokenString string) (string, error) {
	client, err := service.initClient(c)
	if err != nil {
		return "", err
	}

	token, err := service.verifyIDToken(c, client, tokenString)
	if err != nil {
		return "", err
	}

	tx := service.db.Begin()
	defer initializers.CommitOrRollback(tx)
	err = service.repository.Save(tx,
		&models.User{
			ID:       token.Subject,
			Fullname: token.Claims["name"].(string),
			Picture:  token.Claims["picture"].(string),
			Email:    token.Claims["email"].(string),
		},
	)
	if err != nil {
		return "", err
	}

	return client.SessionCookie(c, tokenString, EXPIRES_IN)
}

func (service *AuthService) SignOut(c *gin.Context) error {
	client, err := service.initClient(c)
	if err != nil {
		return err
	}
	return service.revokeToken(c, client)
}

func (service *AuthService) verifyIDToken(c *gin.Context, client *auth.Client, tokenString string) (*auth.Token, error) {
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

func (service *AuthService) VerifySessionToken(c *gin.Context) (*auth.Token, error) {
	client, err := service.initClient(c)
	if err != nil {
		return nil, err
	}
	return service.getSessionToken(c, client)
}
