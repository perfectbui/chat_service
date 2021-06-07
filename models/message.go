package models

import (
	"encoding/json"
	"log"

	"github.com/perfectbui/chat/models/enum"
)

type Message struct {
	Action   enum.ActionValue `json:"action"`
	Message  string           `json:"message"`
	RoomName string           `json:"roomName"`
	Sender   int64            `json:"sender"`
}

func (message *Message) Encode() []byte {
	json, err := json.Marshal(message)
	if err != nil {
		log.Println(err)
	}
	return json
}

func Decode(message []byte) *Message {
	var buf *Message
	err := json.Unmarshal(message, &buf)
	if err != nil {
		log.Println(err)
		return nil
	}
	return buf
}
