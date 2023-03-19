package loginsrv

import (
	"hero_stroy/biz_server/base"
	"hero_stroy/biz_server/mod/user/userdao"
	"hero_stroy/biz_server/mod/user/userdata"
	"hero_stroy/comm/async_op"
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
		user := userdao.GetUserByName(userName)
		nowTime := time.Now().UnixMilli()
		if user == nil {
			user = &userdata.User{
				UserName:   userName,
				Password:   password,
				HeroAvatar: "Hero_Hammer",
				CreateTime: nowTime,
			}
		}
		user.LastLoginTime = nowTime
		err := userdao.SaveOrUpdate(user)
		if err != nil {
			return
		}
		bizResult.SetReturnObj(user)
	}

	continueWith := bizResult.DoComplete
	async_op.Process(async_op.StrToBindId(userName),
		&asyncOp,
		&continueWith)

	return bizResult
}
