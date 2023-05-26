package service

import (
	"errors"
	"strings"
	"time"

	authRepository "github.com/alkuinvito/sirkelin/app/auth/repository"
	"github.com/alkuinvito/sirkelin/initializers"

	"firebase.google.com/go/auth"
	"github.com/gin-gonic/gin"
)

const EXPIRES_IN = time.Hour * 24

type AuthService struct {
	ctx        *gin.Context
	client     *auth.Client
	token      *auth.Token
	repository *authRepository.AuthRepository
	err        error
}

type IAuthService interface {
	Error() error
	getTokenFromCtx() *AuthService
	Init(c *gin.Context) *AuthService
	revokeToken() *AuthService
	setError(err error) *AuthService
	setToken(token *auth.Token) *AuthService
	SignIn(tokenString string) *AuthService
	SignOut() *AuthService
	Token() *auth.Token
	verifyIDToken(tokenString string) *AuthService
}

func (svc *AuthService) Init(c *gin.Context) *AuthService {
	client, err := initializers.InitializeAppDefault().Auth(c)
	authRepo := authRepository.Init()

	return &AuthService{
		ctx:        c,
		client:     client,
		repository: authRepo,
		err:        err,
	}
}

func (svc *AuthService) Error() error {
	return svc.err
}

func (svc *AuthService) getTokenFromCtx() *AuthService {
	bearerToken := svc.ctx.GetHeader("Authorization")
	tokenString := strings.Split(bearerToken, " ")

	if len(tokenString) != 2 {
		cookies, err := svc.ctx.Cookie("session")
		if err != nil {
			return svc.setError(err)
		}
		token, err := svc.client.VerifySessionCookieAndCheckRevoked(svc.ctx, cookies)
		return svc.setToken(token).setError(err)
	}

	token, err := svc.client.VerifySessionCookieAndCheckRevoked(svc.ctx, tokenString[0])
	return svc.setToken(token).setError(err)
}

func (svc *AuthService) revokeToken() *AuthService {
	svc.getTokenFromCtx()
	if svc.Error() != nil {
		return svc
	}

	if err := svc.client.RevokeRefreshTokens(svc.ctx, svc.token.UID); err != nil {
		return svc.setError(err)
	}

	return svc
}

func (svc *AuthService) setError(err error) *AuthService {
	svc.err = err
	return svc
}

func (svc *AuthService) setToken(token *auth.Token) *AuthService {
	svc.token = token
	return svc
}

func (svc *AuthService) SignIn(tokenString string) (string, error) {
	svc.verifyIDToken(tokenString)
	if svc.Error() != nil {
		return "", svc.Error()
	}

	err := svc.repository.AuthenticateByIDToken(
		svc.token.Subject,
		svc.token.Claims["name"].(string),
		svc.token.Claims["picture"].(string),
		svc.token.Claims["email"].(string),
	)
	if err != nil {
		return "", svc.Error()
	}

	return svc.client.SessionCookie(svc.ctx, tokenString, EXPIRES_IN)
}

func (svc *AuthService) SignOut() error {
	return svc.revokeToken().Error()
}

func (svc *AuthService) Token() *auth.Token {
	return svc.getTokenFromCtx().token
}

func (svc *AuthService) verifyIDToken(tokenString string) *AuthService {
	token, err := svc.client.VerifyIDTokenAndCheckRevoked(svc.ctx, tokenString)
	if err != nil {
		if err.Error() == "ID token has been revoked" {
			return svc.setError(errors.New("user must reauthenticate"))
		} else {
			return svc.setError(errors.New("invalid id token"))
		}
	}

	if time.Now().Unix()-token.IssuedAt > 5*60 {
		return svc.setError(errors.New("recent sign-in required"))
	}

	return svc.setToken(token).setError(err)
}
