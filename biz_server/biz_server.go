package main

import (
	"fmt"
	"github.com/gorilla/websocket"
	myWebsocket "hero_stroy/biz_server/network/websocket"
	"hero_stroy/comm/log"
	"net/http"
	"os"
	"path"
)

var upgrader = &websocket.Upgrader{
	ReadBufferSize:  2048,
	WriteBufferSize: 2048,
	//跨域
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func main() {
	fmt.Println("start bizServer")
	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}

	log.Config(path.Dir(ex) + "/log/biz_server.log")
	log.Info("bizServer start success")

	http.HandleFunc("/websocket", websocketHandShake)

	_ = http.ListenAndServe(":12345", nil)

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
	ctx := &myWebsocket.CmdContextImpl{
		Conn: conn,
	}
	go ctx.LoopSendMsg()

	ctx.LoopReadMsg()
}
