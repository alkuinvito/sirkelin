package utils

import (
	"errors"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
)

func CreateToken(uid string) (string, error) {
	jti := uuid.New().String()
	exp := jwt.NewNumericDate(time.Now().Add(time.Hour))
	claims := jwt.RegisteredClaims{
		Subject:   uid,
		ID:        jti,
		ExpiresAt: exp,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, err := token.SignedString([]byte(os.Getenv("SECRET_KEY")))
	if err != nil {
		return "", err
	}
	return ss, nil
}

func ValidateToken(c *gin.Context) error {
	claims, err := DecodeTokenClaims(c)
	if err != nil {
		return err
	}

	if !claims.VerifyExpiresAt(time.Now().Unix(), true) {
		return errors.New("token expired")
	}

	return nil
}

func DecodeTokenClaims(c *gin.Context) (jwt.MapClaims, error) {
	token, err := ExtractTokenHeader(c)
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("failed mapping claims")
	}

	return claims, nil
}

func ExtractTokenHeader(c *gin.Context) (*jwt.Token, error) {
	bearerToken := c.GetHeader("Authorization")
	tokenString := strings.Split(bearerToken, " ")

	return jwt.Parse(tokenString[1], func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(os.Getenv("SECRET_KEY")), nil
	})
}

func GetTokenSubject(c *gin.Context) (string, error) {
	claims, err := DecodeTokenClaims(c)
	if err != nil {
		return "", err
	}

	return claims["sub"].(string), nil
}