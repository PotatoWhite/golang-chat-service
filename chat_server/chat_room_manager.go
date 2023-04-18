package chat_server

import (
	"context"
	"fmt"
	"github.com/go-redis/redis"
	"github.com/google/uuid"
	"log"
	"sync"
	"time"
)

// singleton
var server *ChatServer
var mu sync.Mutex

type ChatServer struct {
	Rooms sync.Map
	ctx   context.Context
	rc    *redis.Client
}

func InitChatServer(ctx context.Context, rc *redis.Client) *ChatServer {
	//mu.Lock()
	//defer mu.Unlock()

	// if server is not initialized, initialize it
	if server == nil {
		server = &ChatServer{
			Rooms: sync.Map{},
			ctx:   ctx,
			rc:    rc,
		}
	}

	return server
}

func (s *ChatServer) OpenRoom(title string, owner *User) *chatroom {
	room := New(title, owner)
	s.Rooms.Store(room.Id, room)
	room.AddUser(owner)

	go server.StartRoomChannel(room.Id)
	return room
}


// send close message to all users in chatroom and delete chatroom from server
func (s *ChatServer) CloseRoom(roomId uuid.UUID, user *User) error {
	value, ok := s.Rooms.Load(roomId)
	if !ok {
		return fmt.Errorf("chatroom %s does not exist", roomId)
	}

	room := value.(*chatroom)
	if room != nil {
		// if user is owner, close chatroom
		if room.Owner.Id == user.Id {
			room.Close()
			s.Rooms.Delete(roomId)
			log.Printf("Room %s is closed", room.Id)
		} else {
			return fmt.Errorf("user %s is not owner of chatroom %s", user.Id, roomId)
		}
	} else {
		return fmt.Errorf("chatroom %s does not exist", roomId)
	}
	return nil
}

func (s *ChatServer) JoinRoom(roomId uuid.UUID, user *User) error {
	value, ok := s.Rooms.Load(roomId)
	if !ok {
		return fmt.Errorf("chatroom %s does not exist", roomId)
	}

	room := value.(*chatroom)

	// if chatroom exists, add user to chatroom
	if room != nil {
		room.AddUser(user)
	} else {
		return fmt.Errorf("chatroom %s does not exist", roomId)
	}

	return nil
}

func (s *ChatServer) LeaveRoom(roomId uuid.UUID, user *User) error {
	value, ok := s.Rooms.Load(roomId)
	if !ok {
		return fmt.Errorf("chatroom %s does not exist", roomId)
	}

	room := value.(*chatroom)
	if room != nil {
		room.RemoveUser(user)
	} else {
		return fmt.Errorf("chatroom %s does not exist", roomId)
	}

	return nil
}

func (s *ChatServer) GetRoom(roomId uuid.UUID) *chatroom {
	value, ok := s.Rooms.Load(roomId)
	if !ok {
		return nil
	}

	return value.(*chatroom)
}

// StartMonitor monitor and logging rooms and close rooms that are empty during the duration
func (s *ChatServer) StartMonitor(second time.Duration) {
	go func() {
		for {
			select {
			case <-s.ctx.Done():
				return
			default:
				s.Rooms.Range(func(key, value interface{}) bool {
					room := value.(*chatroom)
					if len(room.users) == 0 {
						s.CloseRoom(room.Id, room.Owner)
						defer room.Close()
						s.Rooms.Delete(key)
						log.Printf("Room %s is closed", room.Id)
					} else {
						log.Printf("Room %s has %d users", room.Id, len(room.users))
					}
					return true
				})
				time.Sleep(second * time.Second)
			}
		}
	}()
}
