package handler

import (
	"google.golang.org/protobuf/types/dynamicpb"
	"hero_story/biz_server/base"
	"hero_story/biz_server/mod/user/userdata"
	"hero_story/biz_server/msg"
	"hero_story/comm/log"
)

func init() {
	cmdHandlerMap[uint16(msg.MsgCode_WHO_ELSE_IS_HERE_CMD.Number())] = whoElseIsHereCmdHandler
}

func whoElseIsHereCmdHandler(ctx base.MyCmdContext, _ *dynamicpb.Message) {
	if nil == ctx || ctx.GetUserId() <= 0 {
		return
	}
	log.Info("收到还有谁消息, userId: %d", ctx.GetUserId())
	result := &msg.WhoElseIsHereResult{}

	userAll := userdata.GetUserGroup().GetUserAll()
	for _, user := range userAll {
		if user == nil {
			continue
		}
		userInfo := &msg.WhoElseIsHereResult_UserInfo{
			UserId:     uint32(user.UserId),
			UserName:   user.UserName,
			HeroAvatar: user.HeroAvatar,
		}
		if nil != user.MoveState {
			userInfo.MoveState = &msg.WhoElseIsHereResult_UserInfo_MoveState{
				FromPosX:  user.MoveState.FromPosX,
				FromPosY:  user.MoveState.FromPosY,
				ToPosX:    user.MoveState.ToPosX,
				ToPosY:    user.MoveState.ToPosY,
				StartTime: user.MoveState.StartTime,
			}
		}
		result.UserInfo = append(result.UserInfo, userInfo)

	}
	ctx.Write(result)

}
