package controllers

import (
	"net/http"

	"github.com/alkuinvito/sirkelin/models"
	"github.com/alkuinvito/sirkelin/utils"
	"github.com/gin-gonic/gin"
)

func CreateRoom(c *gin.Context) {
	var req models.CreateRoomParams
	var err error

	err = c.ShouldBindJSON(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"data": gin.H{
				"error": "invalid create room request body",
			},
		})
	}

	uid, err := utils.NewAuth(c).GetSession().GetUserID()
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"data": gin.H{
				"error": "invalid bearer token",
			},
		})
		return
	}

	req.Users = append(req.Users, &models.User{ID: uid})

	id, err := models.InsertRoom(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
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

func SendMessage(c *gin.Context) {
	var req models.SendMessageParams
	var param models.RoomIDParams
	var err error

	err = c.ShouldBindUri(&param)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"data": gin.H{
				"error": "invalid room id",
			},
		})
		return
	}

	uid, err := utils.NewAuth(c).GetSession().GetUserID()
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"data": gin.H{
				"error": "invalid bearer token",
			},
		})
		return
	}

	err = c.ShouldBindJSON(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"data": gin.H{
				"error": "invalid message body",
			},
		})
		return
	}

	message := models.Message{
		Body:   req.Body,
		UserID: uid,
		RoomID: param.RoomID,
	}

	err = models.InsertMessage(&message)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"data": gin.H{
				"error": "failed to store message",
			},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"id": message.ID,
		},
	})
}

func GetMessages(c *gin.Context) {
	var err error
	var req models.RoomIDParams
	var room models.Room

	err = c.ShouldBindUri(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"data": gin.H{
				"error": "Invalid room id",
			},
		})
		return
	}

	room.ID = req.RoomID
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

func GetRoomList(c *gin.Context) {
	uid, err := utils.NewAuth(c).GetSession().GetUserID()
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"data": gin.H{
				"error": "invalid bearer token",
			},
		})
		return
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"data": gin.H{
				"error": "failed to retrieve room list",
			},
		})
		return
	}

	roomList := models.RoomList(uid)

	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"rooms": roomList,
		},
	})
}
