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

func NewFirebaseClient(c *gin.Context) (*auth.Client, error) {
	return initializers.InitializeAppDefault().Auth(c)
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

func GetSessionFromContext(c *gin.Context) (string, error) {
	bearerToken := c.GetHeader("Authorization")
	token := strings.Split(bearerToken, " ")

	if len(token) != 2 {
		cookies, err := c.Cookie("session")
		if err != nil {
			return "", errors.New("no session provided")
		}
		return cookies, nil
	}

	return token[1], nil
}

func GetIDFromSession(client *auth.Client, c *gin.Context, session string) (string, error) {
	token, err := client.VerifySessionCookieAndCheckRevoked(c, session)
	if err != nil {
		return "", err
	}

	return token.UID, nil
}
