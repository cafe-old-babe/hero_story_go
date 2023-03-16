package handler

import (
	"github.com/gorilla/websocket"
	"google.golang.org/protobuf/types/dynamicpb"
	"hero_stroy/biz_server/msg"
)

func init() {
	cmdHandlerMap[uint16(msg.MsgCode_USER_ENTRY_CMD)] = userEntryCmdHandler
}

func userEntryCmdHandler(conn *websocket.Conn, pb *dynamicpb.Message) {

}
