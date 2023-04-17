package chat_message_api

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"study02-chat-service/chat_server"
)

var chatSvr *chat_server.ChatServer

func AddChatMsgApis(grp *gin.RouterGroup, srv *chat_server.ChatServer) error {
	chatSvr = srv

	grp.GET("observer/", handleObserver)
	return nil
}

func handleObserver(context *gin.Context) {
	// upgrade http to websocket
	conn, err := upgrader.Upgrade(context.Writer, context.Request, nil)
	if err != nil {
		http.Error(context.Writer, err.Error(), http.StatusBadRequest)
		return
	}
	defer conn.Close()

	handleMessage := func(message []byte) error {
		// write message to client
		if err := conn.WriteMessage(websocket.TextMessage, message); err != nil {
			http.Error(context.Writer, err.Error(), http.StatusBadRequest)
			return err
		}

		log.Printf("write message to client: %s", message)

		return nil
	}

	// read message from client
	for {
		_, rawMsg, err := conn.ReadMessage()
		if err != nil {
			log.Printf("read rawMsg from client error: %s", err.Error())
			continue
		}

		// unmashal rawMsg
		var msg *Message
		if err := json.Unmarshal(rawMsg, &msg); err != nil {
			log.Printf("Unmarshal raw message error: %s", err.Error())
			continue
		}

		switch msg.Type {
		case "join":
			// join room
			room := chatSvr.GetRoom(msg.RoomId)
			if room == nil {
				log.Printf("room not found: %s", msg.RoomId)
				break
			}

			// create observer and register handleMessage
			userIdUUID, err := uuid.Parse(msg.UserId)
			if err != nil {
				log.Printf("userId is invalid %v", msg.UserId)
				break
			}

			room.CreateObserver(userIdUUID, handleMessage)

		case "leave":
			// leave room
			room := chatSvr.GetRoom(msg.RoomId)
			if room == nil {
				log.Printf("room not found: %s", msg.RoomId)
				continue
			}

			// remove observer
			if userId, err := uuid.Parse(msg.UserId); err != nil {
				room.RemoveObserver(userId)
				break
			}
		default:
			break
		}

		log.Printf("read rawMsg from client: %s", rawMsg)
	}
}
