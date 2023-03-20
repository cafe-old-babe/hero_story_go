package handler

import (
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/dynamicpb"
	"hero_story/biz_server/mod/login/loginsrv"
	"hero_story/biz_server/mod/user/userdata"
	"hero_story/biz_server/msg"
	"hero_story/comm/log"
)

func init() {
	cmdHandlerMap[uint16(msg.MsgCode_USER_LOGIN_CMD.Number())] = userLoginCmdHandler
}

func userLoginCmdHandler(ctx MyCmdContext, pbMsgObj *dynamicpb.Message) {
	if ctx == nil || pbMsgObj == nil {
		return
	}
	userLoginCmd := &msg.UserLoginCmd{}
	pbMsgObj.Range(func(f protoreflect.FieldDescriptor, v protoreflect.Value) bool {
		userLoginCmd.ProtoReflect().Set(f, v)
		return true
	})
	log.Info("收到用户消息-->: {%v}", userLoginCmd)

	bizResult := loginsrv.LoginByPasswordAsync(userLoginCmd.GetUserName(), userLoginCmd.GetPassword())
	if bizResult == nil {
		log.Error("业务结果返回空值, userName = %s", userLoginCmd.GetUserName())
		return
	}
	bizResult.OnComplete(func() {
		returnObj := bizResult.GetReturnObj()
		if returnObj == nil {
			log.Error("用户不存在: %s", userLoginCmd.GetUserName())
			return
		}
		user := returnObj.(*userdata.User)

		userLoginResult := &msg.UserLoginResult{
			UserId:     uint32(user.UserId),
			UserName:   user.UserName,
			HeroAvatar: user.HeroAvatar,
		}
		ctx.BindUserId(user.UserId)
		ctx.Write(userLoginResult)
	})

}
