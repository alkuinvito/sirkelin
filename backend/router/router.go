package router

import (
	"os"
	"sirkelin/backend/app/auth/controller"

	"github.com/gin-gonic/gin"
)

type Router struct {
	auth *controller.AuthController
}

func NewRouter(controller *controller.AuthController) *Router {
	return &Router{
		auth: controller,
	}
}

func (router *Router) Handle() *gin.Engine {
	handler := gin.Default()
	gin.SetMode(os.Getenv("APP_MODE"))

	authGroup := handler.Group("/auth")
	{
		authGroup.POST("/sign-in", router.auth.SignIn)
		authGroup.POST("/sign-out", router.auth.SignOut)
	}

	return handler
}
