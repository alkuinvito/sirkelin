package middlewares

import (
	"net/http"

	"github.com/alkuinvito/sirkelin/models"
	"github.com/alkuinvito/sirkelin/utils"
	"github.com/gin-gonic/gin"
)

func RoomAccess() gin.HandlerFunc {
	return func(c *gin.Context) {
		if err := utils.ValidateToken(c); err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"data": gin.H{
					"error": "invalid bearer token",
				},
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

		token, _ := utils.ExtractTokenHeader(c)
		uid, err := utils.GetTokenSubject(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"data": gin.H{
					"error": "invalid bearer token",
				},
			})
			c.Abort()
			return
		}

		room.ID = param.RoomID
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
