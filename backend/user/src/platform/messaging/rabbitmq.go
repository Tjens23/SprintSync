package messaging

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

const (
	EventUserAuthenticated = "user.authenticated"
	EventUserLoggedOut     = "user.logged.out"
	EventUserProfileViewed = "user.profile.viewed"
	EventUserCreated       = "user.created"
	EventUserUpdated       = "user.updated"
	EventUserDeleted       = "user.deleted"
)

type Message struct {
	EventType string                 `json:"event_type"`
	Timestamp time.Time              `json:"timestamp"`
	Data      map[string]interface{} `json:"data"`
}

type MessageHandler func(msg Message) error

type RabbitMQ struct {
	conn         *amqp.Connection
	channel      *amqp.Channel
	exchangeName string
	queueName    string
	handlers     map[string]MessageHandler
}

func NewRabbitMQ(url, exchangeName, queueName string) (*RabbitMQ, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}

	channel, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("failed to open channel: %w", err)
	}

	if err := channel.ExchangeDeclare(
		exchangeName,
		"topic",
		true,
		false,
		false,
		false,
		nil,
	); err != nil {
		channel.Close()
		conn.Close()
		return nil, fmt.Errorf("failed to declare exchange: %w", err)
	}

	queue, err := channel.QueueDeclare(
		queueName,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		channel.Close()
		conn.Close()
		return nil, fmt.Errorf("failed to declare queue: %w", err)
	}

	routingKeys := []string{
		EventUserAuthenticated,
		EventUserLoggedOut,
		EventUserProfileViewed,
		EventUserCreated,
		EventUserUpdated,
		EventUserDeleted,
	}

	for _, key := range routingKeys {
		if err := channel.QueueBind(
			queue.Name,
			key,
			exchangeName,
			false,
			nil,
		); err != nil {
			channel.Close()
			conn.Close()
			return nil, fmt.Errorf("failed to bind queue: %w", err)
		}
	}

	rmq := &RabbitMQ{
		conn:         conn,
		channel:      channel,
		exchangeName: exchangeName,
		queueName:    queueName,
		handlers:     make(map[string]MessageHandler),
	}

	rmq.RegisterHandler(EventUserAuthenticated, rmq.handleUserAuthenticated)
	rmq.RegisterHandler(EventUserLoggedOut, rmq.handleUserLoggedOut)
	rmq.RegisterHandler(EventUserProfileViewed, rmq.handleUserProfileViewed)

	return rmq, nil
}

func (r *RabbitMQ) RegisterHandler(eventType string, handler MessageHandler) {
	r.handlers[eventType] = handler
}

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

	if err := r.channel.Publish(
		r.exchangeName,
		eventType,
		false,
		false,
		amqp.Publishing{
			ContentType:  "application/json",
			Body:         body,
			DeliveryMode: amqp.Persistent,
		},
	); err != nil {
		return fmt.Errorf("failed to publish message: %w", err)
	}

	log.Printf("Published message: %s", eventType)
	return nil
}

func (r *RabbitMQ) StartConsuming() error {
	msgs, err := r.channel.Consume(
		r.queueName,
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("failed to register consumer: %w", err)
	}

	log.Printf("Started consuming messages from queue: %s", r.queueName)
	for msg := range msgs {
		var message Message
		if err := json.Unmarshal(msg.Body, &message); err != nil {
			log.Printf("Error unmarshaling message: %v", err)
			msg.Nack(false, false)
			continue
		}

		if handler, ok := r.handlers[message.EventType]; ok {
			if err := handler(message); err != nil {
				log.Printf("Error handling message %s: %v", message.EventType, err)
				msg.Nack(false, true)
				continue
			}
		} else {
			log.Printf("No handler registered for event type: %s", message.EventType)
		}

		msg.Ack(false)
	}

	return nil
}

func (r *RabbitMQ) handleUserAuthenticated(msg Message) error {
	log.Printf("User authenticated: %+v", msg.Data)
	return nil
}

func (r *RabbitMQ) handleUserLoggedOut(msg Message) error {
	log.Printf("User logged out: %+v", msg.Data)
	return nil
}

func (r *RabbitMQ) handleUserProfileViewed(msg Message) error {
	log.Printf("User profile viewed: %+v", msg.Data)
	return nil
}

func (r *RabbitMQ) Close() {
	if r.channel != nil {
		r.channel.Close()
	}
	if r.conn != nil {
		r.conn.Close()
	}
	log.Println("RabbitMQ connection closed")
}
