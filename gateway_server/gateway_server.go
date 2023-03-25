package main

import (
	"fmt"
	"github.com/gorilla/websocket"
	"hero_story/comm/log"
	gatewaySocket "hero_story/gateway_server/network/websocket"
	"net/http"
	"os"
	"path"
)

var sessionId int32 = 0

var upgrader = &websocket.Upgrader{
	ReadBufferSize:  2048,
	WriteBufferSize: 2048,
	//跨域
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func main() {
	fmt.Println("start gatewayServer")
	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}

	log.Config(path.Dir(ex) + "/log/gateway_server.log")
	log.Info("gateway_server start success")

	http.HandleFunc("/websocket", websocketHandShake)

	_ = http.ListenAndServe(":54321", nil)

}

func websocketHandShake(w http.ResponseWriter, r *http.Request) {
	if nil == w || nil == r {
		return
	}
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Error("websocket upgrade error: %v", err)
	}
	defer func() {
		if err := conn.Close(); err != nil {
			log.Error("websocket close error: %v", err)
		}
	}()
	log.Info("有新客户端加入")
	sessionId += 1
	ctx := &gatewaySocket.CmdContextImpl{
		Conn:      conn,
		SessionId: sessionId,
	}
	ctx.LoopSendMsg()
	ctx.LoopReadMsg()

}
