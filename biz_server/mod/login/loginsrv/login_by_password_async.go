package loginsrv

import (
	"hero_stroy/biz_server/mod/user/userdao"
	"hero_stroy/biz_server/mod/user/userdata"
	"time"
)

//根据用户名称和密码获取用户信息
func LoginByPasswordAsync(userName string, password string) *userdata.User {

	if len(userName) <= 0 || len(password) <= 0 {
		return nil
	}
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
		return nil
	}
	return user
}
