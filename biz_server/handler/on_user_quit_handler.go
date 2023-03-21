package handler

import (
	"hero_story/biz_server/base"
	"hero_story/biz_server/msg"
	"hero_story/biz_server/network/broadcaster"
	"hero_story/comm/log"
)

func OnUserQuitHandler(ctx base.MyCmdContext) {

	if nil == ctx || ctx.GetUserId() <= 0 {
		return
	}
	log.Info("用户离线：%d", ctx.GetUserId())

	// 发送离线消息
	broadcaster.Broadcast(&msg.UserQuitResult{QuitUserId: uint32(ctx.GetUserId())})

}
