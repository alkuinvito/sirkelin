package controllers

import (
	"net/http"

	"github.com/alkuinvito/malakh-api/middlewares"
	"github.com/alkuinvito/malakh-api/models"
	"github.com/alkuinvito/malakh-api/utils"
	"github.com/gin-gonic/gin"
)

func RoomHandler(rg *gin.RouterGroup) {
	rg.GET("/", getRoomList)

	create := rg.Group("/create")
	create.POST("/", createRoom)

	messageHandler := rg.Group("/:id")
	messageHandler.Use(middlewares.RoomPrivillege())
	messageHandler.POST("/", sendMessage)
	messageHandler.GET("/", getMessages)
}

func createRoom(c *gin.Context) {
	var newRoom models.Room

	c.Bind(&newRoom)

	uid, err := utils.ExtractTokenUser(c)

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"data": gin.H{
				"error": "Invalid token",
			},
		})
	}

	newRoom.Users = append(newRoom.Users, &models.User{ID: uid})

	id, err := models.InsertRoom(&newRoom)
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

type MessageBody struct {
	UserID uint
	Body   string
}

func sendMessage(c *gin.Context) {
	var roomId models.RoomId
	var body MessageBody

	if err := c.ShouldBindUri(&roomId); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"data": gin.H{
				"error": "Invalid room id",
			},
		})
		return
	}

	uid, err := utils.ExtractTokenUser(c)

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"data": gin.H{
				"error": "Invalid token",
			},
		})
	}

	c.Bind(&body)

	message := models.Message{
		Body:   body.Body,
		UserID: uid,
		RoomID: roomId.ID,
	}

	err = models.InsertMessage(&message)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"data": gin.H{
				"error": "Failed to store message",
			},
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"data": gin.H{
				"id": message.ID,
			},
		})
	}
}

func getMessages(c *gin.Context) {
	var err error
	var roomId models.RoomId
	var room models.Room

	if err = c.ShouldBindUri(&roomId); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"data": gin.H{
				"error": "Invalid room id",
			},
		})
		return
	}

	room.ID = roomId.ID
	if err = room.PullMessages(); err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"data": gin.H{
				"error": "Room invalid and/or not exist",
			},
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"messages": room.Messages,
		},
	})
}

func getRoomList(c *gin.Context) {
	id, err := utils.ExtractTokenUser(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"data": gin.H{
				"error": "Failed to retrieve room list",
			},
		})
		return
	}

	roomList := models.RoomList(id)

	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"rooms": roomList,
		},
	})
}
