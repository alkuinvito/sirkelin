package controller

import (
	"net/http"
	authService "sirkelin/backend/app/auth/service"
	roomService "sirkelin/backend/app/room/service"
	"sirkelin/backend/models"

	"github.com/gin-gonic/gin"
)

type RoomController struct {
	authService *authService.AuthService
	roomService *roomService.RoomService
}

type IRoomController interface {
	CreateRoom(c*gin.Context)
}

func NewRoomController(authService *authService.AuthService, roomService *roomService.RoomService) *RoomController {
	return &RoomController{
		authService: authService,
		roomService: roomService,
	}
}

func (controller *RoomController) CreateRoom(c *gin.Context) {
	var req models.CreateRoomParams
	var err error

	err = c.ShouldBindJSON(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid create room request body",
		})
	}

	token, err := controller.authService.VerifySessionToken(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "invalid bearer token",
		})
		return
	}

	req.Users = append(req.Users, &models.User{ID: token.UID})

	roomID, err := controller.roomService.Create(req.Users)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create new room",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"id": roomID,
		},
	})
}