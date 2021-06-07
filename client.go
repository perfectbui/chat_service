// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/perfectbui/chat/config"
	"github.com/perfectbui/chat/job"
	"github.com/perfectbui/chat/models"
	"github.com/perfectbui/chat/models/enum"

	"github.com/gorilla/websocket"
	"github.com/streadway/amqp"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// Client is a middleman between the websocket connection and the hub.
type Client struct {
	wsServer *WsServer

	UserID int64 `json:"userID"`
	// The websocket connection.
	conn *websocket.Conn

	// Buffered channel of outbound messages.
	send chan []byte

	rooms map[*Room]bool
}

var ctx = context.Background()

func (c *Client) readPump() {
	defer func() {
		c.conn.Close()
	}()
	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error { c.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}
		decodedMessage := models.Decode(message)
		if decodedMessage == nil {
			return
		}
		decodedMessage.Sender = c.UserID
		c.handleMessage(decodedMessage)
	}
}

func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()
	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The hub closed the channel.
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			// Add queued chat messages to the current websocket message.
			n := len(c.send)
			for i := 0; i < n; i++ {
				w.Write(newline)
				w.Write(<-c.send)
			}

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

func (c *Client) joinRoom(name string) {
	room := c.wsServer.findRoomByName(name)
	if room == nil {
		newRoom := c.wsServer.createRoom(name)
		if newRoom != nil {
			c.rooms[newRoom] = true
			newRoom.register <- c
		}
	} else {
		c.rooms[room] = true
		room.register <- c
	}
}

func (c *Client) leaveRoom(name string) {
	room := c.wsServer.findRoomByName(name)
	if room != nil {
		delete(c.rooms, room)
	}
}

func (c *Client) sendMessage(message *models.Message) {
	for room := range c.wsServer.rooms {
		if room.getName() == message.RoomName {
			go job.AddMessageJob(message.Encode())
			room.publishRoomMessage(message.Encode())
			// room.broadcast <- message
			// room.produceMessage(message.encode())
		}
	}
}

func (room *Room) consumeMessage() {
	ch, err := config.ConsumeConnect.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	err = ch.ExchangeDeclare(
		"logs",
		"fanout",
		true,
		false,
		false,
		false,
		nil,
	)
	failOnError(err, "Failed to declare an exchange")

	q, err := ch.QueueDeclare(
		"",    // name
		false, // durable
		false, // delete when unused
		true,  // exclusive
		false, // no-wait
		nil,   // arguments
	)
	failOnError(err, "Failed to declare a queue")

	err = ch.QueueBind(
		q.Name, // queue name
		"",     // routing key
		"logs", // exchange
		false,
		nil,
	)
	failOnError(err, "Failed to bind a queue")

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	failOnError(err, "Failed to register a consumer")

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			log.Printf(" [x] %s", d.Body)
			room.broadcast <- d.Body
		}
	}()

	log.Printf(" [*] Waiting for logs. To exit press CTRL+C")
	<-forever
}

func (room *Room) produceMessage(body []byte) {

	ch, err := config.ProduceConn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	err = ch.ExchangeDeclare(
		"logs",   // name
		"fanout", // type
		true,     // durable
		false,    // auto-deleted
		false,    // internal
		false,    // no-wait
		nil,      // arguments
	)
	failOnError(err, "Failed to declare an exchange")

	err = ch.Publish(
		"logs", // exchange
		"",     // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        body,
		})

	failOnError(err, "Failed to publish a message")

	log.Printf(" [x] Sent %s", body)

}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func (room *Room) publishRoomMessage(message []byte) {
	err := config.Redis.Publish(ctx, room.getName(), message).Err()

	if err != nil {
		log.Println(err)
	}
}

func (room *Room) subscribeToRoomMessages() {
	pubsub := config.Redis.Subscribe(ctx, room.getName())

	ch := pubsub.Channel()

	for msg := range ch {
		room.broadcast <- []byte(msg.Payload)
	}
}

func (c *Client) handleMessage(decodedMessage *models.Message) {
	switch decodedMessage.Action {
	case enum.Action.JOIN_ROOM:
		c.joinRoom(decodedMessage.RoomName)
	case enum.Action.LEAVE_ROME:
		c.leaveRoom(decodedMessage.RoomName)
	case enum.Action.SEND_MESSAGE:
		c.sendMessage(decodedMessage)
	}
}

func serveWs(wsServer *WsServer, w http.ResponseWriter, r *http.Request, userID int64) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	client := &Client{wsServer: wsServer, conn: conn, UserID: userID, send: make(chan []byte), rooms: make(map[*Room]bool)}
	go client.writePump()
	go client.readPump()
}
