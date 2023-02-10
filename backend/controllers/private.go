package controllers

import (
	"net/http"

	"github.com/alkuinvito/malakh-api/models"
	"github.com/alkuinvito/malakh-api/utils"
	"github.com/gin-gonic/gin"
)

func PrivateHandler(rg *gin.RouterGroup) {
	rg.GET("/", getPrivateList)

	rg.POST("/create", createPrivateRoom)
}

func getPrivateList(c *gin.Context) {
	id, err := utils.ExtractTokenUser(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"data": gin.H{
				"error": "Failed to retrieve room list",
			},
		})
		return
	}

	roomList := models.PrivateList(id)

	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"rooms": roomList,
		},
	})
}

func createPrivateRoom(c *gin.Context) {
	var privateRoom models.Room

	c.Bind(&privateRoom)

	privateRoom.IsPrivate = true

	uid, err := utils.ExtractTokenUser(c)

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