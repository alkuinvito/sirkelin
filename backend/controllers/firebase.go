package controllers

import (
	"fmt"
	"log"
	"net/http"

	firebase "firebase.google.com/go"
	"github.com/alkuinvito/malakh-api/initializers"
	"github.com/gin-gonic/gin"
)

var app = initializers.InitializeAppDefault()

func FirebaseHandler(rg *gin.RouterGroup) {
	rg.POST("/verify", verifyIDToken)
}

type IDToken struct {
	IDToken string
}

func verifyIDToken(ctx *gin.Context) {
	var idToken IDToken
	ctx.Bind(&idToken)
	err := verifyIDTokenAndCheckRevoked(ctx, app, idToken.IDToken)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"data": gin.H{
				"error": err.Error(),
			},
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"data": "",
	})
}

func verifyIDTokenAndCheckRevoked(ctx *gin.Context, app *firebase.App, idToken string) error {
	client, err := app.Auth(ctx)
	if err != nil {
		return err
	}
	token, err := client.VerifyIDTokenAndCheckRevoked(ctx, idToken)
	if err != nil {
		if err.Error() == "ID token has been revoked" {
			return fmt.Errorf("user must reauthenticate session")
		} else {
			return fmt.Errorf("id token is invalid")
		}
	}
	log.Printf("Verified ID token: %v\n", token)

	return nil
}
