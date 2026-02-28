package main

import (
	"SprintSync/src/platform/authenticator"
	"SprintSync/src/platform/messaging"
	"SprintSync/src/platform/router"
	"log"
	"os"

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

	// Initialize RabbitMQ connection
	rabbitMQ, err := messaging.NewRabbitMQ(
		os.Getenv("RABBITMQ_URL"),
		os.Getenv("RABBITMQ_EXCHANGE"),
	)
	if err != nil {
		log.Fatalf("Failed to initialize RabbitMQ: %v", err)
	}
	defer rabbitMQ.Close()

	app := router.New(auth, rabbitMQ)

	log.Print("Server listening on http://localhost:3001/")
	app.Listen(":3001")
}
