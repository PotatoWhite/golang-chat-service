package chat_api

type OpenRoomRequest struct {
	Title        string
	UserNickname string
}

type JoinRoomRequest struct {
	RoomId       string
	UserNickname string
}
