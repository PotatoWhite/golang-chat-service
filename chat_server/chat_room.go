package chat_server

import (
	"github.com/google/uuid"
	"sync"
)

type chatroom struct {
	Id    uuid.UUID
	Title string
	ch    chan string
	Owner *User
	users map[string]*User
	mu    sync.Mutex
}

func (r *chatroom) Error() string {
	//TODO implement me
	panic("implement me")
}

func (r *chatroom) AddUser(user *User) {
	r.users[user.Id] = user
}

func (r *chatroom) RemoveUser(user *User) {
	delete(r.users, user.Id)
}

func (r *chatroom) Broadcast(messageType int, message []byte) {
	r.ch <- string(message)
}

func New(title string, owner *User) *chatroom {

	room := &chatroom{
		Id:    uuid.New(),
		ch:    make(chan string),
		Title: title,
		Owner: owner,
		users: make(map[string]*User),
	}

	return room
}

func (r *chatroom) Close() {
	close(r.ch)
}
