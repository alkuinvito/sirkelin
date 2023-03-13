package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/alkuinvito/malakh-api/controllers"
	"github.com/alkuinvito/malakh-api/initializers"
	"github.com/alkuinvito/malakh-api/middlewares"
	"github.com/gin-gonic/gin"
)

func init() {
	initializers.LoadEnvVar()
	initializers.ConnectToDB()
}

func main() {
	router := gin.Default()

	firebaseGroup := router.Group("/firebase")
	controllers.FirebaseHandler(firebaseGroup)

	privateGroup := router.Group("/private")
	privateGroup.Use(middlewares.RoomAccess())
	controllers.PrivateHandler(privateGroup)

	roomGroup := router.Group("/room")
	roomGroup.Use(middlewares.RoomAccess())
	controllers.RoomHandler(roomGroup)

	userGroup := router.Group("/user")
	controllers.UserHandler(userGroup)

	srv := &http.Server{
		Addr:    os.Getenv("PORT"),
		Handler: router,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}

	select {
	case <-ctx.Done():
		log.Println("timeout of 5 seconds.")
	}
	log.Println("Server exiting")
}
