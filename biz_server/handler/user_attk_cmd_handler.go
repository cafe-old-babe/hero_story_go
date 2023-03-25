package handler

import (
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/dynamicpb"
	"hero_story/biz_server/base"
	"hero_story/biz_server/mod/user/user_data"
	"hero_story/biz_server/mod/user/user_lso"
	"hero_story/biz_server/msg"
	"hero_story/biz_server/network/broadcaster"
	"hero_story/comm/lazy_save"
)

func init() {
	cmdHandlerMap[uint16(msg.MsgCode_USER_ATTK_CMD.Number())] = userAttkCmdHandler
}

func userAttkCmdHandler(ctx base.MyCmdContext, dp *dynamicpb.Message) {
	if ctx == nil || ctx.GetUserId() <= 0 || dp == nil {
		return
	}

	userAttkCmd := &msg.UserAttkCmd{}
	dp.Range(func(descriptor protoreflect.FieldDescriptor, value protoreflect.Value) bool {
		userAttkCmd.ProtoReflect().Set(descriptor, value)
		return true
	})
	userAttkResult := &msg.UserAttkResult{
		AttkUserId:   uint32(ctx.GetUserId()),
		TargetUserId: userAttkCmd.TargetUserId,
	}
	broadcaster.Broadcast(userAttkResult)

	user := user_data.GetUserGroup().GetByUserId(int64(userAttkCmd.TargetUserId))
	if nil == user {
		return
	}

	var subtractHp int32 = 10
	user.CurrHp -= subtractHp
	userSubtractHpResult := &msg.UserSubtractHpResult{
		TargetUserId: userAttkCmd.TargetUserId,
		SubtractHp:   uint32(subtractHp),
	}
	broadcaster.Broadcast(userSubtractHpResult)

	userLso := user_lso.GetUserLso(user)

	lazy_save.SaveOrUpdate(userLso)

}
