package user_dao

import (
	"hero_story/biz_server/base"
	"hero_story/biz_server/mod/user/user_data"
	"hero_story/comm/log"
)

const sqlGetUserByName = `SELECT user_id,user_name,password,hero_avatar, curr_hp FROM t_user WHERE user_name =?`

// GetUserByName 根据用户名查询用户信息
func GetUserByName(userName string) *user_data.User {
	if len(userName) == 0 {
		return nil
	}
	row := base.MysqlDB.QueryRow(sqlGetUserByName, userName)
	if nil == row {
		return nil
	}
	user := &user_data.User{}
	if err := row.Scan(&user.UserId, &user.UserName, &user.Password, &user.HeroAvatar, &user.CurrHp); err != nil {
		log.Error("GetUserByName error: %v", err)
		return nil
	}
	return user
}
