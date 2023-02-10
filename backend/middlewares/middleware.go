package middlewares

import (
	"net/http"

	"github.com/alkuinvito/malakh-api/models"
	"github.com/alkuinvito/malakh-api/utils"
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

func RoomPrivillege() gin.HandlerFunc {
	return func(c *gin.Context) {
		var roomId models.RoomId
		var room models.Room
		var err error
		var id uint

		if err = c.ShouldBindUri(&roomId); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"data": gin.H{
					"error": "Invalid room id",
				},
			})
			return
		}

		id, err = utils.ExtractTokenUser(c)
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
		if room.GetRoomPrivillege(id) {
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
