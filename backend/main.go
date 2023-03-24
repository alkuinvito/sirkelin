package main

import (
	"context"
	"github.com/alkuinvito/sirkelin/router"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/alkuinvito/sirkelin/initializers"
)

func init() {
	initializers.LoadEnvVar()
	initializers.ConnectToDB()
	initializers.InitRedis()
}

func main() {
	routesHandler := router.Handle()

	srv := &http.Server{
		Addr:    os.Getenv("APP_PORT"),
		Handler: routesHandler,
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
