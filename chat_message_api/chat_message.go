package chat_message_api

import (
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

// chat_message_api.go has two functions:
// 1. JoinRoom
// 2. RelayMessage from redis subs to Websocket Client from browser

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type Message struct {
	RoomId  uuid.UUID
	UserId  string
	Type    string
	Payload string
}
