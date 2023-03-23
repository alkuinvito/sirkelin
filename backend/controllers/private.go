package controllers

import (
	"net/http"

	"github.com/alkuinvito/sirkelin/models"
	"github.com/alkuinvito/sirkelin/utils"
	"github.com/gin-gonic/gin"
)

func GetPrivateList(c *gin.Context) {
	uid, err := utils.GetTokenSubject(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"data": gin.H{
				"error": err.Error(),
			},
		})
		return
	}

	roomList := models.PrivateList(uid)

	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"rooms": roomList,
		},
	})
}

func CreatePrivateRoom(c *gin.Context) {
	var privateRoom models.Room

	c.Bind(&privateRoom)
	privateRoom.IsPrivate = true

	uid, err := utils.GetTokenSubject(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"data": gin.H{
				"error": "Invalid token",
			},
		})
	}

	privateRoom.Users = append(privateRoom.Users, &models.User{ID: uid})

	id, err := models.InsertRoom(&privateRoom)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"data": gin.H{
				"error": "Failed to create new room",
			},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"id": id,
		},
	})
}
