package producer

import (
	"fmt"
	"os"

	amqp "github.com/rabbitmq/amqp091-go"
)

var BrokerCh *amqp.Channel

func StartConnect() {
	// Define RabbitMQ server URL.
	amqpServerURL := os.Getenv("AMQP_SERVER_URL")
	fmt.Println("amqpServerURL ", amqpServerURL)

	// Skip RabbitMQ connection if URL is not provided
	if amqpServerURL == "" {
		fmt.Println("AMQP_SERVER_URL not set, skipping RabbitMQ connection")
		return
	}

	connectRabbitMQ := GetRabbitConnection(amqpServerURL)
	if connectRabbitMQ == nil {
		fmt.Println("Failed to connect to RabbitMQ, continuing without message queue")
		return
	}

	defer connectRabbitMQ.Close()

	// Let's start by opening a channel to our RabbitMQ
	// instance over the connection we have already
	// established.
	channelRabbitMQ, err := connectRabbitMQ.Channel()
	if err != nil {
		fmt.Printf("Failed to open RabbitMQ channel: %v\n", err)
		return
	}
	defer channelRabbitMQ.Close()

	// With the instance and declare Queues that we can
	// publish and subscribe to.
	_, err = channelRabbitMQ.QueueDeclare(
		"ProductsDashboard", // queue name
		true,                // durable
		false,               // auto delete
		false,               // exclusive
		false,               // no wait
		nil,                 // arguments
	)
	if err != nil {
		fmt.Printf("Failed to declare queue: %v\n", err)
		return
	}

	fmt.Println("RabbitMQ connection established successfully")
}

func GetRabbitConnection(amqpServerURL string) *amqp.Connection {
	// Create a new RabbitMQ connection.
	connectRabbitMQ, err := amqp.Dial(amqpServerURL)
	if err != nil {
		fmt.Println("error happened ", err)
		return nil
	}

	return connectRabbitMQ
}

func GetChannel() (*amqp.Channel, error) {
	amqpServerURL := os.Getenv("AMQP_SERVER_URL")
	connectRabbitMQ := GetRabbitConnection(amqpServerURL)

	channelRabbitMQ, err := connectRabbitMQ.Channel()
	if err != nil {
		return nil, err
	}

	return channelRabbitMQ, nil
}
