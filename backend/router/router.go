package router

import (
	"os"
	authController "sirkelin/backend/app/auth/controller"
	roomController "sirkelin/backend/app/room/controller"
	"sirkelin/backend/middlewares"

	"github.com/gin-gonic/gin"
)

type Router struct {
	middleware *middlewares.Middleware
	auth       *authController.AuthController
	room       *roomController.RoomController
}

func NewRouter(middleware *middlewares.Middleware, auth *authController.AuthController, room *roomController.RoomController) *Router {
	return &Router{
		middleware: middleware,
		auth:       auth,
		room:       room,
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

	roomGroup := handler.Group("/room")
	{
		roomGroup.Use(router.middleware.RoomAccess())
		roomGroup.POST("/", router.room.CreateRoom)
	}

	return handler
}
