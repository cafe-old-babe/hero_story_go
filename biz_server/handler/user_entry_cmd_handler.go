package handler

import (
	"google.golang.org/protobuf/types/dynamicpb"
	"hero_story/biz_server/msg"
)

func init() {
	cmdHandlerMap[uint16(msg.MsgCode_USER_ENTRY_CMD)] = userEntryCmdHandler
}

func userEntryCmdHandler(ctx MyCmdContext, pb *dynamicpb.Message) {

}
