package utils

import (
	b64 "encoding/base64"
	"errors"
	"os"
	"strings"

	"firebase.google.com/go/auth"
	"github.com/alkuinvito/sirkelin/initializers"
	"github.com/gin-gonic/gin"
)

const (
	NextJS string = "CLIENT_ID_NEXT"
)

type Auth struct {
	ctx *gin.Context
	client *auth.Client
	session *auth.Token
	err error
}

func NewAuth(c *gin.Context) *Auth {
	client, err := initializers.InitializeAppDefault().Auth(c)
	return &Auth{ctx: c, client: client, err: err}
}

func (a *Auth) Client() *auth.Client {
	return a.client
}

func (a *Auth) Error() error {
	return a.err
}

func (a *Auth) GetSession() (*Auth) {
	bearerToken := a.ctx.GetHeader("Authorization")
	token := strings.Split(bearerToken, " ")

	if len(token) != 2 {
		cookies, err := a.ctx.Cookie("session")
		if err != nil {
			a.err = errors.New("no session provided")
			return a
		}
		a.session, a.err = a.client.VerifySessionCookieAndCheckRevoked(a.ctx, cookies)
		return a
	}

	a.session, a.err = a.client.VerifySessionCookieAndCheckRevoked(a.ctx, token[0])
	return a
}

func (a *Auth) GetUserID() (string, error) {
	if a.err != nil {
		return "", a.err
	}

	return a.session.UID, nil
}

func VerifyClientID(id string) (string, error) {
	decoded, err := b64.StdEncoding.DecodeString(id)
	if err != nil {
		return "", err
	}

	if string(decoded) == os.Getenv(NextJS) {
		return NextJS, nil
	}
	return "", errors.New("unknown client type")
}
