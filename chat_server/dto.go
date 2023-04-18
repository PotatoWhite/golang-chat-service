package chat_server

import "github.com/google/uuid"

type Message struct {
	RoomId       uuid.UUID
	UserNickname string
	Message      string
}

type User struct {
	Id            uuid.UUID                  `json:"Id"`
	Nickname      string                     `json:"nickname"`
	HandleMessage func(message []byte) error `json:"-"`
	LastSeen      int64                      `json:"lastSeen"`
}
