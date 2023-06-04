package middlewares

import (
	"github.com/gin-gonic/gin"
	"net/http"
	authService "sirkelin/backend/app/auth/service"
	roomService "sirkelin/backend/app/room/service"
)

type Middleware struct {
	Engine      *gin.Engine
	authService authService.AuthService
	roomService roomService.RoomService
}

type IMiddleware interface {
	RoomAccess() gin.HandlerFunc
}

func NewMiddleware(engine *gin.Engine, authService authService.AuthService, roomService roomService.RoomService) *Middleware {
	return &Middleware{
		Engine:      engine,
		authService: authService,
		roomService: roomService,
	}
}

func (middleware *Middleware) RoomAccess() gin.HandlerFunc {
	return func(c *gin.Context) {
		_, err := middleware.authService.VerifySessionToken(c)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "invalid session token",
			})
			c.Abort()
			return
		}
		c.Next()
	}
}
