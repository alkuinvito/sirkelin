package middlewares

import (
	"github.com/gin-gonic/gin"
	"net/http"
	authService "sirkelin/backend/app/auth/service"
	roomService "sirkelin/backend/app/room/service"

	"github.com/gin-gonic/gin"
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

func (middleware *Middleware) RoomPrivilege() gin.HandlerFunc {
	return func(c *gin.Context) {
		var param models.RoomIDParams
		var room models.Room
		var err error

		err = c.ShouldBindUri(&param)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"data": gin.H{
					"error": "invalid room id",
				},
			})
			c.Abort()
			return
		}

		token, err := middleware.authService.VerifySessionToken(c)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"data": gin.H{
					"error": "invalid session token",
				},
			})
			c.Abort()
			return
		}

		room.ID = param.RoomID
		uid := token.UID
		if room.GetRoomPrivilege(uid) {
			c.JSON(http.StatusForbidden, gin.H{
				"data": gin.H{
					"error": "user is not member of the room",
				},
			})
			c.Abort()
			return
		}

		c.Next()
	}
}
