package handler

import (
	"hero_story/biz_server/base"
	"hero_story/biz_server/mod/user/userdata"
	"hero_story/biz_server/mod/user/userlso"
	"hero_story/biz_server/msg"
	"hero_story/biz_server/network/broadcaster"
	"hero_story/comm/lazy_save"
	"hero_story/comm/log"
)

func OnUserQuitHandler(ctx base.MyCmdContext) {

	if nil == ctx || ctx.GetUserId() <= 0 {
		return
	}
	log.Info("用户离线：%d", ctx.GetUserId())

	// 发送离线消息
	broadcaster.Broadcast(&msg.UserQuitResult{QuitUserId: uint32(ctx.GetUserId())})

	user := userdata.GetUserGroup().GetByUserId(ctx.GetUserId())

	if user == nil {
		return
	}
	userLso := userlso.GetUserLso(user)
	lazy_save.Discard(userLso)
	log.Info("用户离线, userId: %d, 立即保存", ctx.GetUserId())
	userLso.SaveOrUpdate()

}
