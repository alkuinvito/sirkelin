package initializers

import (
	"context"
	"log"

	firebase "firebase.google.com/go"
)

func InitializeAppDefault() *firebase.App {
	app, err := firebase.NewApp(context.Background(), nil)
	if err != nil {
		log.Fatalf("error initializing app: %v\n", err)
	}

	return app
}
