package handler

import (
	"github.com/gorilla/websocket"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/dynamicpb"
	"hero_stroy/biz_server/msg"
	"hero_stroy/comm/log"
)

func init() {
	cmdHandlerMap[uint16(msg.MsgCode_USER_LOGIN_CMD.Number())] = userLoginCmdHandler
}

func userLoginCmdHandler(conn *websocket.Conn, pbMsgObj *dynamicpb.Message) {
	if conn == nil || pbMsgObj == nil {
		return
	}
	userLoginCmd := &msg.UserLoginCmd{}
	pbMsgObj.Range(func(f protoreflect.FieldDescriptor, v protoreflect.Value) bool {
		userLoginCmd.ProtoReflect().Set(f, v)
		return true
	})
	log.Info("收到用户消息-->: {%v}", userLoginCmd)
}
