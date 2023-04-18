package chat_server

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"log"
	"sync"
	"time"
)

type chatroom struct {
	Id     uuid.UUID
	Title  string
	Owner  *User
	users  map[uuid.UUID]*User
	mu     sync.Mutex         `json:"-"`
	cancel context.CancelFunc `json:"-"`
}

func (r *chatroom) Error() string {
	//TODO implement me
	panic("implement me")
}

func (r *chatroom) AddUser(user *User) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.users[user.Id] = user
}

func (r *chatroom) RemoveUser(user *User) {
	r.mu.Lock()
	defer r.mu.Unlock()

	//cleanup handler
	user.HandleMessage = nil

	//remove user from chatroom
	delete(r.users, user.Id)
}

func New(title string, owner *User) *chatroom {
	_, cancelFunc := context.WithCancel(context.Background())

	room := &chatroom{
		Id:     uuid.New(),
		Title:  title,
		Owner:  owner,
		users:  make(map[uuid.UUID]*User),
		cancel: cancelFunc,
	}

	return room
}

func (r *chatroom) Close() {
	// send closing message to all users
	r.mu.Lock()
	defer r.mu.Unlock()

	for _, user := range r.users {
		user.HandleMessage([]byte("chatroom is closing"))
		user.HandleMessage = nil
	}

	// remove all users
	r.users = make(map[uuid.UUID]*User)

	r.cancel()
}

func (r *chatroom) CreateObserver(userId uuid.UUID, msgHandler func(message []byte) error) error {
	// get user
	user, ok := r.users[userId]
	if !ok {
		log.Printf("user %s does not exist in chatroom %s", userId, r.Id)
		return fmt.Errorf("user %s does not exist in chatroom %s", userId, r.Id)
	}

	// create & set handler
	user.HandleMessage = msgHandler

	// update last seen
	r.UpdateLastSeen(user)

	return nil
}

func (r *chatroom) RemoveObserver(id uuid.UUID) {
	// get user
	user, ok := r.users[id]
	if !ok {
		log.Printf("user %s does not exist in chatroom %s", id, r.Id)
		return
	}

	// remove user
	r.RemoveUser(user)

}

func (r *chatroom) UpdateObserver(id uuid.UUID) {
	// get user
	user, ok := r.users[id]
	if !ok {
		log.Printf("user %s does not exist in chatroom %s", id, r.Id)
		return
	}

	// update user
	r.UpdateLastSeen(user)

}

func (r *chatroom) UpdateLastSeen(user *User) {
	user.LastSeen = time.Now().Unix()
}
