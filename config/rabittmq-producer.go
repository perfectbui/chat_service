package config

import (
	"github.com/streadway/amqp"
)

var ProduceConn *amqp.Connection

func InitProducer() {
	// Define RabbitMQ server URL.
	amqpServerURL := "amqp://localhost:5672"

	// Create a new RabbitMQ connection.
	var err error
	ProduceConn, err = amqp.Dial(amqpServerURL)
	if err != nil {
		panic(err)
	}
	// defer connectRabbitMQ.Close()

	// // Let's start by opening a channel to our RabbitMQ
	// // instance over the connection we have already
	// // established.
	// channelRabbitMQ, err := connectRabbitMQ.Channel()
	// if err != nil {
	// 	panic(err)
	// }
	// defer channelRabbitMQ.Close()

	// // With the instance and declare Queues that we can
	// // publish and subscribe to.
	// _, err = channelRabbitMQ.QueueDeclare(
	// 	"QueueService1", // queue name
	// 	true,            // durable
	// 	false,           // auto delete
	// 	false,           // exclusive
	// 	false,           // no wait
	// 	nil,             // arguments
	// )
	// if err != nil {
	// 	panic(err)
	// }

	// // Create a new Fiber instance.
	// app := fiber.New()

	// // Add middleware.
	// app.Use(
	// 	logger.New(), // add simple logger
	// )

	// // Add route.
	// app.Get("/send", func(c *fiber.Ctx) error {
	// 	// Create a message to publish.
	// 	message := amqp.Publishing{
	// 		ContentType: "text/plain",
	// 		Body:        []byte(c.Query("msg")),
	// 	}

	// 	// Attempt to publish a message to the queue.
	// 	if err := channelRabbitMQ.Publish(
	// 		"",              // exchange
	// 		"QueueService1", // queue name
	// 		false,           // mandatory
	// 		false,           // immediate
	// 		message,         // message to publish
	// 	); err != nil {
	// 		return err
	// 	}

	// 	return nil
	// })

	// // Start Fiber API server.
	// log.Fatal(app.Listen(":3000"))
}
