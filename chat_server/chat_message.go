package chat_server

type PubSubMessage struct {
	messageType int
	message     Message
}

type Message struct {
	RoomId       string
	UserNickname string
	Message      string
}
