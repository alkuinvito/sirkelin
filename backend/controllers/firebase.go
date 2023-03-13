package controllers

import (
	"fmt"
	"log"
	"net/http"

	firebase "firebase.google.com/go"
	"github.com/gin-gonic/gin"
)

func FirebaseHandler(rg *gin.RouterGroup) {
	rg.POST("/verify", verifyIDToken)
}

type IDToken struct {
	Token string `json:"idToken"`
}

func verifyIDToken(ctx *gin.Context) {
	var idToken IDToken
	ctx.Bind(&idToken)
	err := verifyIDTokenAndCheckRevoked(ctx, &firebase.App{}, idToken.Token)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"data": gin.H{
				"error": err,
			},
		})
	}

	ctx.JSON(http.StatusOK, gin.H{
		"data": "",
	})
}

func verifyIDTokenAndCheckRevoked(ctx *gin.Context, app *firebase.App, idToken string) error {
	client, err := app.Auth(ctx)
	if err != nil {
		return fmt.Errorf("failed getting Auth client")
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
