package chat_server

import "github.com/google/uuid"

type PubSubMessage struct {
	messageType int
	message     Message
}

type Message struct {
	RoomId       uuid.UUID
	UserNickname string
	Message      string
}
