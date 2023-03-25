package handler

import (
	"fmt"
	"hero_story/biz_server/base"
	"hero_story/biz_server/mod/user/user_data"
	"hero_story/biz_server/mod/user/user_lock"
	"hero_story/biz_server/mod/user/user_lso"
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

	//用户登录使用的是用户名,用户退出使用的是userId,会分配到两个线程里面,使用登出锁
	//加锁
	lockKey := fmt.Sprintf("UserQuit_%d", ctx.GetUserId())
	user_lock.TryLock(lockKey)

	// 发送离线消息
	broadcaster.Broadcast(&msg.UserQuitResult{QuitUserId: uint32(ctx.GetUserId())})

	user := user_data.GetUserGroup().GetByUserId(ctx.GetUserId())

	if user == nil {
		return
	}
	userLso := user_lso.GetUserLso(user)
	lazy_save.Discard(userLso)
	log.Info("用户离线, userId: %d, 立即保存", ctx.GetUserId())

	userLso.SaveOrUpdate(func() {
		//解锁
		user_lock.UnLock(lockKey)
	})

}
