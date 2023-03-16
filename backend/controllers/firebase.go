package controllers

import (
	"net/http"
	"time"

	"github.com/alkuinvito/malakh-api/initializers"
	"github.com/alkuinvito/malakh-api/models"
	"github.com/alkuinvito/malakh-api/utils"
	"github.com/gin-gonic/gin"
)

var app = initializers.InitializeAppDefault()

func FirebaseHandler(rg *gin.RouterGroup) {
	rg.POST("/", verifyIDToken)
}

type IDToken struct {
	IDToken string
}

func verifyIDToken(c *gin.Context) {
	var idToken IDToken
	c.Bind(&idToken)

	client, err := app.Auth(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"data": gin.H{
				"error": "firebase admin sdk error",
			},
		})
		return
	}

	token, err := client.VerifyIDTokenAndCheckRevoked(c, idToken.IDToken)
	if err != nil {
		if err.Error() == "ID token has been revoked" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"data": gin.H{
					"error": "user must reauthenticate",
				},
			})
		} else {
			c.JSON(http.StatusBadRequest, gin.H{
				"data": gin.H{
					"error": "id token invalid",
				},
			})
		}
		return
	}

	if time.Now().Unix()-token.IssuedAt > 5*60 {
		c.JSON(http.StatusUnauthorized, gin.H{
			"data": gin.H{
				"error": "recent sign-in required",
			},
		})
		return
	}

	user := &models.User{
		ID:       token.Subject,
		Fullname: token.Claims["name"].(string),
		Picture:  token.Claims["picture"].(string),
		Email:    token.Claims["email"].(string),
	}
	user.UserAuthenticate()

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
