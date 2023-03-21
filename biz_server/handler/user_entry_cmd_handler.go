package handler

import (
	"google.golang.org/protobuf/types/dynamicpb"
	"hero_story/biz_server/base"
	"hero_story/biz_server/mod/user/userdata"
	"hero_story/biz_server/msg"
	"hero_story/biz_server/network/broadcaster"
	"hero_story/comm/log"
)

func init() {
	cmdHandlerMap[uint16(msg.MsgCode_USER_ENTRY_CMD)] = userEntryCmdHandler
}

func userEntryCmdHandler(ctx base.MyCmdContext, _ *dynamicpb.Message) {
	if ctx == nil || ctx.GetUserId() <= 0 {
		return
	}

	log.Info("收到用户入场消息, userId: %d", ctx.GetUserId())
	user := userdata.GetUserGroup().GetByUserId(ctx.GetUserId())
	if nil == user {
		log.Error("用户不存在, userId: %d", ctx.GetUserId())
		return
	}
	userEntryResult := &msg.UserEntryResult{
		UserId:     uint32(ctx.GetUserId()),
		UserName:   user.UserName,
		HeroAvatar: user.HeroAvatar,
	}
	broadcaster.Broadcast(userEntryResult)
}
