package chat_server

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
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

	go server.RedisSubscribe(room.Id)
	return room
}

func (s *ChatServer) CloseRoom(roomId uuid.UUID, user *User) error {
	value, ok := s.Rooms.Load(roomId)
	if !ok {
		return fmt.Errorf("chatroom %s does not exist", roomId)
	}

	room := value.(*chatroom)

	// if requesting from Owner, last user, or chatroom is empty, close chatroom
	if room.Owner.Id == user.Id || len(room.users) <= 1 || room == nil {
		room.Broadcast(websocket.CloseMessage, &Message{
			UserNickname: "Operator",
			RoomId:       roomId,
			Message:      fmt.Sprintf("Chatroom %s is closed", room.Title),
		})
		defer room.Close()
	} else {
		// reject request
		return fmt.Errorf("user %s is not the Owner of chatroom %s", user.Id, roomId)
	}

	// delete chatroom from server
	s.Rooms.Delete(roomId)

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

func (s *ChatServer) Broadcast(roomId uuid.UUID, msg Message) error {
	// if chatroom exists, remove user from chatroom
	value, ok := s.Rooms.Load(roomId)
	if !ok {
		return fmt.Errorf("chatroom %s does not exist", roomId)
	}

	room := value.(*chatroom)

	// if chatroom exists, broadcast message to chatroom
	if room != nil {
		mu.Lock()
		s.rc.Publish(roomId.String(), msg)
		mu.Unlock()
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

// RedisSubscribe is listening to redis channel and broadcast message to chatroom
func (s *ChatServer) RedisSubscribe(roomId uuid.UUID) {
	pubsub := s.rc.Subscribe(roomId.String())
	defer pubsub.Close()

	for {
		select {
		case <-s.ctx.Done():
			return
		default:
			msg, err := pubsub.ReceiveMessage()
			if err != nil {
				log.Println(err)
				continue
			}

			var message Message
			err = json.Unmarshal([]byte(msg.Payload), &message)
			if err != nil {
				log.Println(err)
				continue
			}

			room := s.GetRoom(message.RoomId)
			if room != nil {
				for _, user := range room.users {
					var msgBytes []byte
					msgBytes, err = json.Marshal(message)
					if err != nil {
						log.Println(err)
						continue
					}

					if user.HandleMessage != nil {
						user.HandleMessage(msgBytes)
					} else {
						log.Printf("User %s is not connected", user.Id)
					}
				}
			}

		}
	}

}

// PublishMessage is publishing message to redis
func (s *ChatServer) PublishMessage(roomId uuid.UUID, msg *Message) {
	jsonPayload, err := json.Marshal(msg)
	if err != nil {
		log.Println(err)
		return
	}

	publish := s.rc.Publish(roomId.String(), jsonPayload)
	if publish.Err() != nil {
		log.Println(publish.Err())
	}
}
