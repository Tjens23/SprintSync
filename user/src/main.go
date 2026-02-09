package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"SprintSync/src/platform/messaging"
	"SprintSync/src/platform/router"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Failed to load the env vars: %v", err)
	}

	rabbitMQ, err := messaging.NewRabbitMQ(
		os.Getenv("RABBITMQ_URL"),
		os.Getenv("RABBITMQ_EXCHANGE"),
		os.Getenv("RABBITMQ_QUEUE"),
	)
	if err != nil {
		log.Fatalf("Failed to initialize RabbitMQ: %v", err)
	}
	defer rabbitMQ.Close()

	go func() {
		if err := rabbitMQ.StartConsuming(); err != nil {
			log.Printf("Error consuming messages: %v", err)
		}
	}()

	app := router.New(rabbitMQ)

	port := os.Getenv("PORT")
	if port == "" {
		port = "3002"
	}

	go func() {
		if err := app.Listen(":" + port); err != nil {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	log.Printf("User service listening on http://localhost:%s/", port)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	<-sigChan

	log.Println("Shutting down gracefully...")
	if err := app.Shutdown(); err != nil {
		log.Printf("Error during shutdown: %v", err)
	}
}
