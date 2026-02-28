package messaging

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

// Event types
const (
	EventUserAuthenticated = "user.authenticated"
	EventUserLoggedOut     = "user.logged.out"
	EventUserProfileViewed = "user.profile.viewed"
)

// Message represents a message structure
type Message struct {
	EventType string                 `json:"event_type"`
	Timestamp time.Time              `json:"timestamp"`
	Data      map[string]interface{} `json:"data"`
}

// RabbitMQ handles RabbitMQ connections and operations
type RabbitMQ struct {
	conn         *amqp.Connection
	channel      *amqp.Channel
	exchangeName string
}

// NewRabbitMQ creates a new RabbitMQ connection for publishing
func NewRabbitMQ(url, exchangeName string) (*RabbitMQ, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}

	channel, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("failed to open channel: %w", err)
	}

	// Declare exchange
	err = channel.ExchangeDeclare(
		exchangeName, // name
		"topic",      // type
		true,         // durable
		false,        // auto-deleted
		false,        // internal
		false,        // no-wait
		nil,          // arguments
	)
	if err != nil {
		channel.Close()
		conn.Close()
		return nil, fmt.Errorf("failed to declare exchange: %w", err)
	}

	return &RabbitMQ{
		conn:         conn,
		channel:      channel,
		exchangeName: exchangeName,
	}, nil
}

// Publish publishes a message to the exchange
func (r *RabbitMQ) Publish(eventType string, data map[string]interface{}) error {
	msg := Message{
		EventType: eventType,
		Timestamp: time.Now(),
		Data:      data,
	}

	body, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	err = r.channel.Publish(
		r.exchangeName, // exchange
		eventType,      // routing key
		false,          // mandatory
		false,          // immediate
		amqp.Publishing{
			ContentType:  "application/json",
			Body:         body,
			DeliveryMode: amqp.Persistent,
		},
	)
	if err != nil {
		return fmt.Errorf("failed to publish message: %w", err)
	}

	log.Printf("Published message: %s", eventType)
	return nil
}

// Close closes the RabbitMQ connection
func (r *RabbitMQ) Close() {
	if r.channel != nil {
		r.channel.Close()
	}
	if r.conn != nil {
		r.conn.Close()
	}
	log.Println("RabbitMQ connection closed")
}
