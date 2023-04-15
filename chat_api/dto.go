package chat_api

type OpenRoomRequest struct {
	Title        string
	UserNickname string
}

type JoinRoomRequest struct {
	RoomId       string
	UserNickname string
}

type SendMessageRequest struct {
	RoomId       string
	UserNickname string
	Message      string
}
