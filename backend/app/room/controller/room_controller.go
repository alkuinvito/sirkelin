package controller

import (
	"errors"
	"github.com/jackc/pgx/v5/pgconn"
	"net/http"
	"reflect"
	roomService "sirkelin/backend/app/room/service"
	userService "sirkelin/backend/app/user/service"
	"sirkelin/backend/app/websocket"
	"sirkelin/backend/models"
	"sirkelin/backend/utils"

	"github.com/gin-gonic/gin"
)

type RoomController struct {
	userService *userService.UserService
	roomService *roomService.RoomService
	hub         *websocket.Hub
}

type IRoomController interface {
	Connect(c *gin.Context)
	CreateRoom(c *gin.Context)
	Delete(c *gin.Context)
	GetPrivateRooms(c *gin.Context)
	GetRoomById(c *gin.Context)
	GetRooms(c *gin.Context)
	UpdateRoom(c *gin.Context)
}

func NewRoomController(userService *userService.UserService, roomService *roomService.RoomService, hub *websocket.Hub) *RoomController {
	return &RoomController{
		userService: userService,
		roomService: roomService,
		hub:         hub,
	}
}

func (controller *RoomController) Connect(c *gin.Context) {
	websocket.ServeWs(controller.hub, c.Writer, c.Request)
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

	roomID, err := controller.roomService.Create(&req)
	if err != nil {
		if reflect.TypeOf(err) == reflect.TypeOf(&pgconn.PgError{}) {
			if err.(*pgconn.PgError).Code == "23503" {
				c.JSON(http.StatusBadRequest, gin.H{
					"error": "invalid participants id",
				})
				return
			}
		}
		if errors.Is(err, utils.ErrMinimumParticipant) {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to create room",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": roomID,
	})
}

func (controller *RoomController) Delete(c *gin.Context) {
	var err error
	var param models.RoomIDParams

	err = c.ShouldBindUri(&param)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid room id",
		})
		return
	}

	err = controller.roomService.Delete(param.RoomID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to delete room",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": param.RoomID,
	})
}

func (controller *RoomController) GetPrivateRooms(c *gin.Context) {
	var err error

	token, err := controller.userService.VerifySessionToken(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "invalid bearer token",
		})
		return
	}

	rooms, err := controller.roomService.GetPrivateRooms(token.UID)
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

func (controller *RoomController) GetRoomById(c *gin.Context) {
	var err error
	var req models.RoomIDParams

	err = c.ShouldBindUri(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid room id",
		})
		return
	}

	room, err := controller.roomService.GetRoomById(req.RoomID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to retrieve rooms",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": room,
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

func (controller *RoomController) PushMessage(c *gin.Context) {
	var err error
	var param models.RoomIDParams
	var req models.SendMessageParams

	err = c.ShouldBindUri(&param)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid room id",
		})
		return
	}

	err = c.ShouldBind(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	token, err := controller.userService.VerifySessionToken(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "invalid bearer token",
		})
		return
	}

	messageID, err := controller.roomService.PushMessage(token.UID, param.RoomID, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to push message",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": messageID,
	})
}

func (controller *RoomController) UpdateRoom(c *gin.Context) {
	var err error
	var param models.RoomIDParams
	var req models.UpdateRoomSchema

	err = c.ShouldBindUri(&param)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid room id",
		})
		return
	}

	err = c.ShouldBind(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	err = controller.roomService.UpdateRoom(param.RoomID, req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": param.RoomID,
	})
}
