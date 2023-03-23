package utils

import (
	b64 "encoding/base64"
	"errors"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
)

var TokenType string = "Bearer"
var ExpiresIn int = 600

type Identities struct {
	UserID   string `json:"user_id"`
	Fullname string `json:"fullname"`
	Email    string `json:"email"`
}

type RefreshToken struct {
	jwt.RegisteredClaims
}

type AccessToken struct {
	Identities Identities `json:"identities"`
	jwt.RegisteredClaims
}

func CreateRefreshToken(uid string) (string, error) {
	jti := uuid.New().String()
	exp := jwt.NewNumericDate(time.Now().Add(time.Hour))
	iat := jwt.NewNumericDate(time.Now())
	claims := RefreshToken{
		jwt.RegisteredClaims{
			Subject:   uid,
			ID:        jti,
			ExpiresAt: exp,
			IssuedAt:  iat,
		},
	}

	return SignToken(claims)
}

func CreateToken(uid, fullname, email string) (string, error) {
	jti := uuid.New().String()
	exp := jwt.NewNumericDate(time.Now().Add(time.Minute * 10))
	iat := jwt.NewNumericDate(time.Now())
	claims := AccessToken{
		Identities{
			UserID:   uid,
			Fullname: fullname,
			Email:    email,
		},
		jwt.RegisteredClaims{
			Subject:   uid,
			ID:        jti,
			ExpiresAt: exp,
			IssuedAt:  iat,
		},
	}

	return SignToken(claims)
}

func DecodeToken(tokenString string) (*AccessToken, error) {
	token, err := jwt.ParseWithClaims(tokenString, &AccessToken{}, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		decoded := make([]byte, b64.StdEncoding.EncodedLen(len(os.Getenv("SECRET_KEY"))))
		b64.StdEncoding.Encode(decoded, []byte(os.Getenv("SECRET_KEY")))
		return decoded, nil
	})
	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*AccessToken); ok && token.Valid {
		return claims, nil
	}

	return nil, err
}

func ExtractTokenHeader(c *gin.Context) (string, error) {
	bearerToken := c.GetHeader("Authorization")
	tokenString := strings.Split(bearerToken, " ")

	if len(tokenString) != 2 {
		return "", errors.New("invalid authorization header")
	}

	return tokenString[1], nil
}

func GetTokenSubject(token string) (string, error) {
	decoded, err := DecodeToken(token)
	if err != nil {
		return "", err
	}

	return decoded.Identities.UserID, nil
}

func SignToken(claims jwt.Claims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	decoded := make([]byte, b64.StdEncoding.EncodedLen(len(os.Getenv("SECRET_KEY"))))
	b64.StdEncoding.Encode(decoded, []byte(os.Getenv("SECRET_KEY")))
	ss, err := token.SignedString(decoded)
	if err != nil {
		return "", err
	}
	return ss, nil
}

func ValidateToken(c *gin.Context) error {
	token, err := ExtractTokenHeader(c)
	if err != nil {
		return err
	}

	_, err = DecodeToken(token)
	if err != nil {
		return err
	}

	return nil
}
