package controller

import (
	"net/http"
	"os"
	"sirkelin/backend/app/auth/service"

	"github.com/gin-gonic/gin"
)

type GetSessionParams struct {
	ClientID string `json:"client_id"`
	IDToken  string `json:"id_token"`
}

type AuthController struct {
	service *service.AuthService
}

type IAuthController interface {
	SignIn(c *gin.Context)
	SignOut(c *gin.Context)
}

func NewAuthController(authService *service.AuthService) *AuthController {
	return &AuthController{
		service: authService,
	}
}

func (controller *AuthController) SignIn(c *gin.Context) {
	var req GetSessionParams
	var err error

	err = c.Bind(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request body",
		})
		return
	}

	session, err := controller.service.SignIn(c, req.IDToken)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request body",
		})
		return
	}

	c.SetCookie("session", session, int(service.EXPIRES_IN), "/", os.Getenv("APP_HOST"), true, true)
}

func (controller *AuthController) SignOut(c *gin.Context) {
	if err := controller.service.SignOut(c); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid session token",
		})
		return
	}

	c.SetCookie("session", "", 0, "/", os.Getenv("APP_HOST"), true, true)
}
