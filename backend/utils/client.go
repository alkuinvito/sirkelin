package utils

import (
	b64 "encoding/base64"
	"errors"
	"github.com/gin-gonic/gin"
	"os"
)

const (
	NextJS string = "CLIENT_ID_NEXT"
)

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

func SetRefreshMethod(c *gin.Context, clientType, refreshToken string) {
	switch clientType {
	case NextJS:
		c.SetCookie("refresh_token", refreshToken, 3600, "/", os.Getenv("APP_HOST"), false, false)
	default:
		return
	}
}
