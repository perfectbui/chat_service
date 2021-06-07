package job

import (
	"log"

	"github.com/streadway/amqp"
)

var produceConn *amqp.Connection

const CREATE_ROOM_ROUTE_KEY = "create_room"
const ADD_MESSAGE_ROUTE_KEY = "add_message"
const ROOM_EXCHANGE_NAME = "room"

func LoadProducer() {
	amqpServerURL := "amqp://localhost:5672"
	// Create a new RabbitMQ connection.
	var err error
	produceConn, err = amqp.Dial(amqpServerURL)
	if err != nil {
		panic(err)
	}
}

func CreateRoomJob(roomName string) {
	ch, err := produceConn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	err = ch.ExchangeDeclare(
		ROOM_EXCHANGE_NAME, // name
		"direct",           // type
		true,               // durable
		false,              // auto-deleted
		false,              // internal
		false,              // no-wait
		nil,                // arguments
	)
	failOnError(err, "Failed to declare an exchange")

	err = ch.Publish(
		ROOM_EXCHANGE_NAME,    // exchange
		CREATE_ROOM_ROUTE_KEY, // routing key
		false,                 // mandatory
		false,                 // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(roomName),
		})
	failOnError(err, "Failed to publish a message")

	log.Printf(" [x] Sent %s", roomName)
}

func AddMessageJob(message []byte) {
	ch, err := produceConn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	err = ch.ExchangeDeclare(
		ROOM_EXCHANGE_NAME, // name
		"direct",           // type
		true,               // durable
		false,              // auto-deleted
		false,              // internal
		false,              // no-wait
		nil,                // arguments
	)
	failOnError(err, "Failed to declare an exchange")

	err = ch.Publish(
		ROOM_EXCHANGE_NAME,    // exchange
		ADD_MESSAGE_ROUTE_KEY, // routing key
		false,                 // mandatory
		false,                 // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        message,
		})
	failOnError(err, "Failed to publish a message")

}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}
