package chat_server

type chatMessage struct {
	Author   string `json:"author"`
	Contents string `json:"contents"`
}
