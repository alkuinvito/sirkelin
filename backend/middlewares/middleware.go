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
					"error": "Invalid token",
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
		var roomId models.RoomId
		var room models.Room

		if err := c.ShouldBindUri(&roomId); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"data": gin.H{
					"error": "Invalid room id",
				},
			})
			return
		}

		id, err := utils.GetTokenSubject(c)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"data": gin.H{
					"error": "Invalid token",
				},
			})
			c.Abort()
			return
		}

		room.ID = roomId.ID
		if room.GetRoomPrivilege(id) {
			c.JSON(http.StatusForbidden, gin.H{
				"data": gin.H{
					"error": "User is not member of the room",
				},
			})
			c.Abort()
			return
		}

		c.Next()
	}

}
