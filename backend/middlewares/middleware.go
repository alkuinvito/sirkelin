package middlewares

import (
	"github.com/gin-gonic/gin"
	"net/http"
	roomService "sirkelin/backend/app/room/service"
	userService "sirkelin/backend/app/user/service"
	"sirkelin/backend/models"
)

type Middleware struct {
	userService *userService.UserService
	roomService *roomService.RoomService
}

type IMiddleware interface {
	RoomAccess() gin.HandlerFunc
	RoomPrivilege() gin.HandlerFunc
}

func NewMiddleware(userService *userService.UserService, roomService *roomService.RoomService) *Middleware {
	return &Middleware{
		userService: userService,
		roomService: roomService,
	}
}

func (middleware *Middleware) RoomAccess() gin.HandlerFunc {
	return func(c *gin.Context) {
		_, err := middleware.userService.VerifySessionToken(c)
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
		var err error

		err = c.ShouldBindUri(&param)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "invalid room id",
			})
			c.Abort()
			return
		}

		token, err := middleware.userService.VerifySessionToken(c)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "invalid session token",
			})
			c.Abort()
			return
		}

		roomID := param.RoomID
		uid := token.UID
		isParticipant, err := middleware.roomService.CheckRoomParticipant(roomID, uid)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "unable to check participants in current room",
			})
			c.Abort()
			return
		}

		if !isParticipant {
			c.JSON(http.StatusForbidden, gin.H{
				"error": "user is not a participant of this room",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}
