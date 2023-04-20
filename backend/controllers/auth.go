package controllers

import (
	"errors"
	"net/http"
	"os"
	"time"

	"firebase.google.com/go/auth"
	"github.com/alkuinvito/sirkelin/models"
	"github.com/alkuinvito/sirkelin/utils"

	"github.com/gin-gonic/gin"
)

type GetJWTParams struct {
	ClientID string `json:"client_id"`
	IDToken  string `json:"id_token"`
}

type RefreshJWTParams struct {
	ClientID     string `json:"client_id"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func verifyIDToken(c *gin.Context, client *auth.Client, IDToken string) (*auth.Token, error) {
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
	var req GetJWTParams
	var err error

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

	firebase, err := utils.NewFirebaseClient(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"data": gin.H{
				"error": "firebase admin sdk error",
			},
		})
		return
	}

	token, err := verifyIDToken(c, firebase, req.IDToken)
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

	err = models.AuthenticateByIDToken(token)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"data": gin.H{
				"error": err.Error(),
			},
		})
		return
	}

	switch client {
	case utils.NextJS:
		expiresIn := time.Hour * 24
		session, err := firebase.SessionCookie(c, req.IDToken, expiresIn)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"data": gin.H{
					"error": err.Error(),
				},
			})
			return
		}

		c.SetCookie("session", session, int(expiresIn.Seconds()), "/", os.Getenv("APP_HOST"), true, true)
		c.JSON(http.StatusOK, gin.H{
			"data": gin.H{
				"expires_in": int(expiresIn),
			},
		})
		return
	default:
		expiresIn := time.Hour * 24 * 5
		session, err := firebase.SessionCookie(c, req.IDToken, expiresIn)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"data": gin.H{
					"error": err.Error(),
				},
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"data": gin.H{
				"session_token": session,
				"expires_in":    int(expiresIn),
			},
		})
	}
}

func SignOut(c *gin.Context) {
	var err error

	firebase, err := utils.NewFirebaseClient(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"data": gin.H{
				"error": "firebase admin sdk error",
			},
		})
		return
	}

	session, err := utils.GetSessionFromContext(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"data": gin.H{
				"error": "invalid session token",
			},
		})
		return
	}
	uid, _ := utils.GetIDFromSession(firebase, c, session)

	err = firebase.RevokeRefreshTokens(c, uid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"data": gin.H{
				"error": "firebase admin sdk error",
			},
		})
		return
	}

	c.SetCookie("session", "", 0, "/", os.Getenv("APP_HOST"), true, true)
}