package controllers

import (
	"net/http"

	"github.com/alkuinvito/sirkelin/models"
	"github.com/alkuinvito/sirkelin/utils"
	"github.com/gin-gonic/gin"
)

func GetPrivateList(c *gin.Context) {
	uid, err := utils.NewAuth(c).GetSession().GetUserID()
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"data": gin.H{
				"error": "invalid bearer token",
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
	req.IsPrivate = true

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
				"error": "failed to create new room",
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
