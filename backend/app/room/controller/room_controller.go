package controller

import (
	"net/http"
	roomService "sirkelin/backend/app/room/service"
	userService "sirkelin/backend/app/user/service"
	"sirkelin/backend/models"

	"github.com/gin-gonic/gin"
)

type RoomController struct {
	userService *userService.UserService
	roomService *roomService.RoomService
}

type IRoomController interface {
	CreateRoom(c *gin.Context)
	GetRooms(c *gin.Context)
}

func NewRoomController(userService *userService.UserService, roomService *roomService.RoomService) *RoomController {
	return &RoomController{
		userService: userService,
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

	token, err := controller.userService.VerifySessionToken(c)
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

func (controller *RoomController) GetRooms(c *gin.Context) {
	var err error

	token, err := controller.userService.VerifySessionToken(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "invalid bearer token",
		})
		return
	}

	rooms, err := controller.roomService.GetRooms(token.UID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to retrieve rooms",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": rooms,
	})
}
