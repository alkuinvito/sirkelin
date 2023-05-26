package controller

import (
	authService "github.com/alkuinvito/sirkelin/app/auth/service"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

type GetSessionParams struct {
	ClientID string `json:"client_id"`
	IDToken  string `json:"id_token"`
}

type AuthController struct {
	service *authService.AuthService
}

type IAuthController interface {
	SignIn(c *gin.Context)
	SignOut(c *gin.Context)
}

func Init() *AuthController {
	return &AuthController{}
}

func (ctr *AuthController) SignIn(c *gin.Context) {
	var req GetSessionParams
	var err error

	err = c.Bind(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request body",
		})
		return
	}

	ctr.service.Init(c)
	session, err := ctr.service.SignIn(req.IDToken)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request body",
		})
		return
	}

	c.SetCookie("session", session, int(authService.EXPIRES_IN), "/", os.Getenv("APP_HOST"), true, true)
}

func (ctr *AuthController) SignOut(c *gin.Context) {
	ctr.service.Init(c)
	if ctr.service.Error() != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "firebase admin sdk error",
		})
		return
	}

	if err := ctr.service.SignOut(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid session token",
		})
		return
	}

	c.SetCookie("session", "", 0, "/", os.Getenv("APP_HOST"), true, true)
}
