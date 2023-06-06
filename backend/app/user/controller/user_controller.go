package controller

import (
	"net/http"
	"os"
	"sirkelin/backend/app/user/service"

	"github.com/gin-gonic/gin"
)

type GetSessionParams struct {
	ClientID string `json:"client_id"`
	IDToken  string `json:"id_token"`
}

type UserController struct {
	service *service.UserService
}

type IUserController interface {
	GetAll(c *gin.Context)
	SignIn(c *gin.Context)
	SignOut(c *gin.Context)
}

func NewUserController(userService *service.UserService) *UserController {
	return &UserController{
		service: userService,
	}
}

func (controller *UserController) GetAll(c *gin.Context) {
	users, err := controller.service.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to retrieve users",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"users": users,
		},
	})
}

func (controller *UserController) SignIn(c *gin.Context) {
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

func (controller *UserController) SignOut(c *gin.Context) {
	if err := controller.service.SignOut(c); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid session token",
		})
		return
	}

	c.SetCookie("session", "", 0, "/", os.Getenv("APP_HOST"), true, true)
}
