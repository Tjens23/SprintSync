package main

import (
	"SprintSync/src/platform/authenticator"
	"SprintSync/src/platform/router"
	"log"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Failed to load the env vars: %v", err)
	}

	auth, err := authenticator.New()
	if err != nil {
		log.Fatalf("Failed to initialize the authenticator: %v", err)
	}

	app := router.New(auth)

	log.Print("Server listening on http://localhost:3001/")
	app.Listen(":3001")
}
