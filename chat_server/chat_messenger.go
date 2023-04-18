package chat_server

import (
	"encoding/json"
	"github.com/google/uuid"
	"log"
)

// StartRoomChannel is listening to redis channel and broadcast message to chatroom
func (s *ChatServer) StartRoomChannel(roomId uuid.UUID) {
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
