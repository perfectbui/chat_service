package main

import (
	"fmt"

	"github.com/streadway/amqp"
)

type Room struct {
	name           string
	clients        map[*Client]bool
	register       chan *Client
	unregister     chan *Client
	broadcast      chan []byte
	produceChannel *amqp.Channel
}

func (r *Room) getName() string {
	return r.name
}

func newRoom(name string) *Room {
	// channelRabbitMQ, err := config.ProduceConn.Channel()
	// if err != nil {
	// 	panic(err)
	// }

	// // With the instance and declare Queues that we can
	// // publish and subscribe to.
	// err = channelRabbitMQ.ExchangeDeclare(
	// 	name,     // name
	// 	"fanout", // type
	// 	true,     // durable
	// 	false,    // auto-deleted
	// 	false,    // internal
	// 	false,    // no-wait
	// 	nil,      // arguments
	// )

	// if err != nil {
	// 	panic(err)
	// }
	// return &Room{
	// 	name:           name,
	// 	clients:        make(map[*Client]bool),
	// 	register:       make(chan *Client),
	// 	unregister:     make(chan *Client),
	// 	broadcast:      make(chan []byte),
	// 	produceChannel: channelRabbitMQ,
	// }

	return &Room{
		name:       name,
		clients:    make(map[*Client]bool),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		broadcast:  make(chan []byte),
		// produceChannel: channelRabbitMQ,
	}
}

func (r *Room) runRoom() {
	go r.subscribeToRoomMessages()
	// go r.consumeMessage()
	for {
		select {
		case client := <-r.register:
			fmt.Println("nhan dang ki vao room tu client", client.UserID, "vao room", r.name)
			r.clients[client] = true
		case client := <-r.unregister:
			if _, ok := r.clients[client]; ok {
				delete(r.clients, client)
				close(client.send)
			}
		case message := <-r.broadcast:
			fmt.Println("nhan messsage tu broadcast")
			fmt.Println("so clients", len(r.clients))
			for client := range r.clients {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(r.clients, client)
				}
			}
		}
	}
}
