package main

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"study02-chat-service/chat_api"
	"study02-chat-service/chat_redis"
	"study02-chat-service/chat_server"
	"study02-chat-service/config"
	"study02-chat-service/log"
)

func main() {
	// load config
	cfg := config.LoadConfigOrExit()
	log.Infof("config: %+v", cfg)

	// connect to chat_redis
	rc, e := chat_redis.OpenNewRedisClient(&cfg.Redis)
	if e != nil {
		log.Fatalf("failed to init chat_redis, %v", e)
	} else {
		log.Infof("chat_redis connected : %s:%d", cfg.Redis.Host, cfg.Redis.Port)
		defer rc.Close()
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// start chat server
	svr := chat_server.InitChatServer(ctx, rc)
	svr.StartMonitor(1)

	// init gin
	eng := gin.Default()
	chatGrp := eng.Group("/chat")
	chatApi := chat_api.AddChatApis(chatGrp, svr)
	if chatApi != nil {
		log.Fatalf("failed to add chat apis")
	}

	//msgGrp := eng.Group("/ws")
	//msgApi := chat_message.AddMessageApis(msgGrp, svr)
	//if msgApi != nil {
	//	log.Fatalf("failed to add message apis")
	//}

	// start http server
	if e := eng.Run(fmt.Sprintf(":%d", cfg.Server.Port)); e != nil {
		log.Fatalf("failed to start http server, %v", eng)
	} else {
		log.Infof("http server started : %s:%d", cfg.Server.Host, cfg.Server.Port)
	}
}
