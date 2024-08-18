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

	connectRabbitMQ := GetRabbitConnection(amqpServerURL)

	defer connectRabbitMQ.Close()

	// Let's start by opening a channel to our RabbitMQ
	// instance over the connection we have already
	// established.
	channelRabbitMQ, err := connectRabbitMQ.Channel()
	if err != nil {
		panic(err)
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
		panic(err)
	}
}

func GetRabbitConnection(amqpServerURL string) *amqp.Connection {
	// Create a new RabbitMQ connection.
	connectRabbitMQ, err := amqp.Dial(amqpServerURL)
	if err != nil {
		// panic(err)
		fmt.Println("error happened ", err)
	}

	return connectRabbitMQ
}
