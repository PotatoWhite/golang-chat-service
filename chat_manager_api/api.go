package chat_manager_api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
	"study02-chat-service/chat_server"
)

var chatSvr *chat_server.ChatServer

func AddChatApis(grp *gin.RouterGroup, srv *chat_server.ChatServer) error {
	chatSvr = srv

	grp.GET(":roomId", handleGetRoomInfo)
	grp.POST("", handleCreateAndJoinRoom)
	grp.POST(":roomId/join", handleJoinRoom)
	grp.POST(":roomId/message", handleSendMessage)
	grp.DELETE(":roomId", handleCloseRoom)

	return nil
}

func handleGetRoomInfo(c *gin.Context) {
	roomId, err := uuid.Parse(c.Param("roomId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	room := chatSvr.GetRoom(roomId)
	if room == nil {
		c.JSON(http.StatusNoContent, gin.H{
			"error": "room not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"room": room,
	})
}

func handleCreateAndJoinRoom(c *gin.Context) {
	// get Request Body
	var rb OpenRoomRequest
	if err := c.ShouldBindJSON(&rb); err != nil {
		badRequest(c, err)
		return
	}

	// create room
	room := chatSvr.OpenRoom(rb.Title, &chat_server.User{
		Id:       uuid.New(),
		Nickname: rb.UserNickname,
	})

	c.JSON(http.StatusOK, gin.H{
		"room": room,
	})
}

func handleJoinRoom(c *gin.Context) {
	roomId, err := uuid.Parse(c.Param("roomId"))
	if err != nil {
		badRequest(c, err)
		return
	}

	room := chatSvr.GetRoom(roomId)
	if room == nil {
		noContent(c, nil)
		return
	}

	// get Request Body
	var rb JoinRoomRequest
	if err := c.ShouldBindJSON(&rb); err != nil {
		badRequest(c, err)
		return
	}

	// join room
	room.AddUser(&chat_server.User{
		Id:       uuid.New(),
		Nickname: rb.UserNickname,
	})

	c.JSON(http.StatusOK, gin.H{
		"room": room,
	})
}

func handleCloseRoom(c *gin.Context) {
	roomId, err := uuid.Parse(c.Param("roomId"))
	if err != nil {
		badRequest(c, err)
		return
	}

	// get userId from query
	userId := c.Query("userId")
	if userId == "" {
		badRequest(c, err)
		return
	}

	userIdUuid, err := uuid.Parse(userId)
	if err != nil {
		badRequest(c, err)
		return
	}

	// close room
	room := chatSvr.CloseRoom(roomId, &chat_server.User{
		Id: userIdUuid,
	})

	if room == nil {
		// success
		c.JSON(http.StatusOK, gin.H{
			"room": room,
		})
	}

	// internal server error
	internalServerError(c, fmt.Errorf("failed to close room : %v", room))
}

func handleSendMessage(c *gin.Context) {
	roomId, err := uuid.Parse(c.Param("roomId"))
	if err != nil {
		badRequest(c, err)
		return
	}

	room := chatSvr.GetRoom(roomId)
	if room == nil {
		noContent(c, nil)
		return
	}

	// get Request Body
	var rb SendMessageRequest
	if err := c.ShouldBindJSON(&rb); err != nil {
		badRequest(c, err)
		return
	}

	// send message to Redis
	chatSvr.PublishMessage(roomId, &chat_server.Message{
		RoomId:       roomId,
		UserNickname: rb.UserNickname,
		Message:      rb.Message,
	})

	c.JSON(http.StatusOK, gin.H{
		"room": room,
	})
}
