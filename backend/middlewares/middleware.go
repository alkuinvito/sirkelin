package middlewares

import (
	"net/http"

	authService "github.com/alkuinvito/sirkelin/app/auth/service"
	"github.com/gin-gonic/gin"
)

func RoomAccess() gin.HandlerFunc {
	return func(c *gin.Context) {
		auth := authService.Init(c)
		if auth.Error() != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "invalid session token",
			})
			c.Abort()
			return
		}
		c.Next()
	}
}

func RoomPrivilege() gin.HandlerFunc {
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

		auth := authService.Init(c)
		if auth.Error() != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"data": gin.H{
					"error": "invalid session token",
				},
			})
			c.Abort()
			return
		}

		room.ID = param.RoomID
		uid := auth.Token().UID
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
