package middlewares

import (
	"github.com/gin-gonic/gin"
	"github.com/gorilla/csrf"
	adapter "github.com/gwatts/gin-adapter"
	"log"
	"net/http"
	"os"
	roomService "sirkelin/backend/app/room/service"
	userService "sirkelin/backend/app/user/service"
	"sirkelin/backend/models"
)

type Middleware struct {
	userService *userService.UserService
	roomService *roomService.RoomService
}

type IMiddleware interface {
	AuthenticatedUser() gin.HandlerFunc
	AuthorizedUser() gin.HandlerFunc
	CSRF() gin.HandlerFunc
}

func NewMiddleware(userService *userService.UserService, roomService *roomService.RoomService) *Middleware {
	return &Middleware{
		userService: userService,
		roomService: roomService,
	}
}

func (middleware *Middleware) AuthenticatedUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		_, err := middleware.userService.VerifySessionToken(c)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "invalid session token",
			})
			c.Abort()
			return
		}
		c.Next()
	}
}

func (middleware *Middleware) AuthorizedUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		var param models.RoomIDParams
		var err error

		err = c.ShouldBindUri(&param)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "invalid room id",
			})
			c.Abort()
			return
		}

		token, err := middleware.userService.VerifySessionToken(c)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "invalid session token",
			})
			c.Abort()
			return
		}

		roomID := param.RoomID
		uid := token.UID
		isParticipant, err := middleware.roomService.CheckRoomParticipant(roomID, uid)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "unable to check participants in current room",
			})
			c.Abort()
			return
		}

		if !isParticipant {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "room with this id does not exist",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

func (middleware *Middleware) CSRF() gin.HandlerFunc {
	return adapter.Wrap(
		csrf.Protect(
			[]byte(os.Getenv("CSRF_KEY")),
			csrf.MaxAge(0),
			csrf.Secure(true),
			csrf.ErrorHandler(http.HandlerFunc(
				func(w http.ResponseWriter, r *http.Request) {
					log.Println(r.Header["X-Csrf-Token"])
					w.WriteHeader(http.StatusForbidden)
					w.Write([]byte(`{ "error": "csrf token mismatch" }`))
				}),
			),
		),
	)
}
