package userdao

import (
	"hero_story/biz_server/base"
	"hero_story/biz_server/mod/user/userdata"
	"hero_story/comm/log"
)

const sqlGetUserByName = `SELECT user_id,user_name,password,hero_avatar, curr_hp FROM t_user WHERE user_name =?`

// GetUserByName 根据用户名查询用户信息
func GetUserByName(userName string) *userdata.User {
	if len(userName) == 0 {
		return nil
	}
	row := base.MysqlDB.QueryRow(sqlGetUserByName, userName)
	if nil == row {
		return nil
	}
	user := &userdata.User{}
	if err := row.Scan(&user.UserId, &user.UserName, &user.Password, &user.HeroAvatar, &user.CurrHp); err != nil {
		log.Error("GetUserByName error: %v", err)
		return nil
	}
	return user
}
