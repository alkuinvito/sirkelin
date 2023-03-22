package controllers

import (
	"errors"
	"firebase.google.com/go/auth"
	"net/http"
	"time"

	"github.com/alkuinvito/sirkelin/initializers"
	"github.com/alkuinvito/sirkelin/models"
	"github.com/alkuinvito/sirkelin/utils"
	"github.com/gin-gonic/gin"
)

var app = initializers.InitializeAppDefault()

func AuthHandler(rg *gin.RouterGroup) {
	rg.POST("/sign-in", signIn)
}

type AuthenticationResponse struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
}

func verifyIDToken(c *gin.Context) (*auth.Token, error) {
	var idToken string
	var err error

	err = c.Bind(&idToken)
	if err != nil {
		return nil, errors.New("failed retrieving id token")
	}

	client, err := app.Auth(c)
	if err != nil {
		return nil, errors.New("firebase admin sdk error")
	}

	token, err := client.VerifyIDTokenAndCheckRevoked(c, idToken)
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

	return token, nil
}

func signIn(c *gin.Context) {
	var err error

	token, err := verifyIDToken(c)
	if err != nil {
		if err.Error() == "failed retrieving id token" {
			c.JSON(http.StatusBadRequest, gin.H{
				"data": gin.H{
					"error": err.Error(),
				},
			})
		} else if err.Error() == "firebase admin sdk error" {
			c.JSON(http.StatusInternalServerError, gin.H{
				"data": gin.H{
					"error": err.Error(),
				},
			})
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{
				"data": gin.H{
					"error": err.Error(),
				},
			})
		}
	}

	models.AuthenticateByIDToken(token)

	jwt, err := utils.CreateToken(token.Subject)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"data": gin.H{
				"error": "failed to generate jwt",
			},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"access_token": jwt,
		},
	})
}
