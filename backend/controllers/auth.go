package controllers

import (
	"errors"
	"firebase.google.com/go/auth"
	"github.com/alkuinvito/sirkelin/models"
	"github.com/alkuinvito/sirkelin/utils"
	"net/http"
	"time"

	"github.com/alkuinvito/sirkelin/initializers"
	"github.com/gin-gonic/gin"
)

type SignInRequest struct {
	ClientID string `json:"client_id,required"`
	IDToken  string `json:"id_token,required"`
}

type SignInResponse struct {
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
	AccessToken string `json:"access_token"`
}

func NewFirebaseClient(c *gin.Context) (*auth.Client, error) {
	var app = initializers.InitializeAppDefault()
	return app.Auth(c)
}

func verifyIDToken(c *gin.Context, IDToken string) (*auth.Token, error) {
	var err error

	client, err := NewFirebaseClient(c)
	if err != nil {
		return nil, errors.New("firebase admin sdk error")
	}

	token, err := client.VerifyIDTokenAndCheckRevoked(c, IDToken)
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

func SignIn(c *gin.Context) {
	var err error
	var req SignInRequest

	err = c.Bind(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"data": gin.H{
				"error": "invalid request body",
			},
		})
		return
	}

	client, err := utils.VerifyClientID(req.ClientID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"data": gin.H{
				"error": "unknown client type",
			},
		})
		return
	}

	switch client {
	case utils.NextJS:

	}

	token, err := verifyIDToken(c, req.IDToken)
	if err != nil {
		if err.Error() == "firebase admin sdk error" {
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
		return
	}

	models.AuthenticateByIDToken(token)

	accessToken, err := utils.CreateToken(token.Subject, token.Claims["name"].(string), token.Claims["email"].(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"data": gin.H{
				"error": "failed to generate jwt",
			},
		})
		return
	}

	refreshToken, _ := utils.CreateRefreshToken(token.Subject)
	utils.SetRefreshMethod(c, client, refreshToken)
	c.JSON(http.StatusOK, gin.H{
		"data": SignInResponse{
			TokenType:   utils.TokenType,
			ExpiresIn:   utils.ExpiresIn,
			AccessToken: accessToken,
		},
	})
}

func RefreshTokens(c *gin.Context) {

}
