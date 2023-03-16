package main

import (
	"encoding/binary"
	"fmt"
	"github.com/gorilla/websocket"
	"hero_stroy/biz_server/handler"
	"hero_stroy/biz_server/msg"
	"hero_stroy/comm/log"
	"hero_stroy/comm/main_thread"
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
	/*http.HandleFunc("/websocket", func(w http.ResponseWriter, request *http.Request) {
		_, _ = w.Write([]byte("Hello , World!\n"))
	})*/
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
	for {
		_, msgData, err := conn.ReadMessage()
		if err != nil {
			log.Error("websocket readMessage error: %v ", err)
			break
		}
		log.Info("msgData: %v", msgData)

		/*
			descriptor := msg.File_GameMsgProtocol_proto.Messages().ByName("UserLoginCmd")
			message := dynamicpb.NewMessage(descriptor)
			cmd := &msg.UserLoginCmd{}
			message.Range(func(dp protoreflect.FieldDescriptor, val protoreflect.Value) bool {
				cmd.ProtoReflect().Set(dp, val)
				return true
			})
		*/
		msgCode := binary.BigEndian.Uint16(msgData[2:4])
		message, err := msg.Decode(msgData[4:], int16(msgCode))
		if err != nil {
			log.Error("message message msgCode: %d, err: %v+", msgCode, err)
			continue
		}
		log.Info("收到客户端消息,,msgCode: %d, message Name: %v", msgCode, message.Descriptor().Name())

		cmdHandlerFunc := handler.CreateCmdHandler(msgCode)
		if cmdHandlerFunc == nil {
			log.Error("没有查询到指令处理函数,msgCode: %d", msgCode)
			continue
		}
		main_thread.Process(func() {
			cmdHandlerFunc(conn, message)

		})
		/*


			//region
			descriptor := msg.File_GameMsgProtocol_proto.Messages().ByName("UserLoginCmd")
			name := message.Get(descriptor.Fields().ByName("userName"))

			password := message.Get(descriptor.Fields().ByName("password"))

			log.Info("userName: %s, password: %s", name, password)
			result := &msg.UserLoginResult{
				UserId:     1,
				UserName:   name.String(),
				HeroAvatar: "Hero_Shaman",
			}
			byteArray, err := msg.Encode(result)
			if err != nil {
				log.Error("message Encode result error: %v", err)
				continue
			}

			if err = conn.WriteMessage(websocket.BinaryMessage, byteArray); err != nil {
				log.Error("writeMessage result error: %v", err)
			}
			//endregion
		*/
	}
}
