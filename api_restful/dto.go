package api_restful

type OpenRoomRequest struct {
	Title        string
	UserNickname string
}

type JoinRoomRequest struct {
	RoomId       string
	UserNickname string
}

type SendMessageRequest struct {
	UserNickname string
	Message      string
}
