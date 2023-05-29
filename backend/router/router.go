package router

import (
	authController "github.com/alkuinvito/sirkelin/app/auth/controller"
	"os"

	"github.com/gin-gonic/gin"
)

func Handle() *gin.Engine {
	router := gin.Default()
	gin.SetMode(os.Getenv("APP_MODE"))

	ac := authController.Init()
	authGroup := router.Group("/auth")
	{
		authGroup.POST("/sign-in", ac.SignIn)
		authGroup.POST("/sign-out", ac.SignOut)
	}

	return router
}
