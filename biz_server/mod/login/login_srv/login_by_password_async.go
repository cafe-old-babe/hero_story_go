package login_srv

import (
	"fmt"
	"hero_story/biz_server/base"
	"hero_story/biz_server/mod/user/user_dao"
	"hero_story/biz_server/mod/user/user_data"
	"hero_story/biz_server/mod/user/user_lock"
	"hero_story/comm/async_op"
	"time"
)

// LoginByPasswordAsync 根据用户名称和密码获取用户信息
func LoginByPasswordAsync(userName string, password string) *base.AsyncBizResult {
	//31 * h java compiler (h << 5) - h
	if len(userName) <= 0 || len(password) <= 0 {
		return nil
	}
	bizResult := &base.AsyncBizResult{}
	asyncOp := func() {
		user := user_dao.GetUserByName(userName)
		nowTime := time.Now().UnixMilli()
		if user == nil {
			user = &user_data.User{
				UserName:   userName,
				Password:   password,
				HeroAvatar: "Hero_Hammer",
				CreateTime: nowTime,
			}
		}

		//检查是否有登出锁,如果有锁,直接退出
		lockKey := fmt.Sprintf("UserQuit_%d", user.UserId)
		if user_lock.HashLock(lockKey) {
			bizResult.SetReturnObj(nil)
			return
		}
		user.LastLoginTime = nowTime
		err := user_dao.SaveOrUpdate(user)
		if err != nil {
			return
		}
		bizResult.SetReturnObj(user)
	}

	async_op.Process(async_op.StrToBindId(userName),
		&asyncOp, nil)

	return bizResult
}
