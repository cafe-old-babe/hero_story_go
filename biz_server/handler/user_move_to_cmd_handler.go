package handler

import (
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/dynamicpb"
	"hero_story/biz_server/base"
	"hero_story/biz_server/mod/user/user_data"
	"hero_story/biz_server/msg"
	"hero_story/biz_server/network/broadcaster"
	"time"
)

func init() {
	cmdHandlerMap[uint16(msg.MsgCode_USER_MOVE_TO_CMD.Number())] = userMoveToCmdCmdHandler
}

func userMoveToCmdCmdHandler(ctx base.MyCmdContext, dp *dynamicpb.Message) {
	if ctx == nil || ctx.GetUserId() <= 0 || dp == nil {
		return
	}

	user := user_data.GetUserGroup().GetByUserId(ctx.GetUserId())
	if user == nil {
		return
	}
	cmd := &msg.UserMoveToCmd{}
	dp.Range(func(descriptor protoreflect.FieldDescriptor, value protoreflect.Value) bool {
		cmd.ProtoReflect().Set(descriptor, value)
		return true
	})
	if user.MoveState == nil {
		user.MoveState = &user_data.MoveState{}
	}
	nowTime := uint64(time.Now().UnixMilli())
	user.MoveState.FromPosX = cmd.MoveFromPosX
	user.MoveState.FromPosY = cmd.MoveFromPosY
	user.MoveState.ToPosY = cmd.MoveToPosY
	user.MoveState.ToPosX = cmd.MoveToPosX
	user.MoveState.StartTime = nowTime
	result := &msg.UserMoveToResult{
		MoveUserId:    uint32(ctx.GetUserId()),
		MoveFromPosX:  cmd.MoveFromPosX,
		MoveFromPosY:  cmd.MoveFromPosY,
		MoveToPosX:    cmd.MoveToPosX,
		MoveToPosY:    cmd.MoveToPosY,
		MoveStartTime: nowTime,
	}
	broadcaster.Broadcast(result)
}
