package chat_redis

import (
	"github.com/google/uuid"
	"study02-chat-service/chat_server"
)

type Message struct {
	Author  chat_server.User
	Message string
	RoomId  uuid.UUID
}
