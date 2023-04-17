package chat_server

import "github.com/google/uuid"

type User struct {
	Id            uuid.UUID                  `json:"Id"`
	Nickname      string                     `json:"nickname"`
	HandleMessage func(message []byte) error `json:"-"`
}
