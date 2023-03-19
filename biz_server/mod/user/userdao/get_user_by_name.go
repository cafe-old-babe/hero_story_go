package userdao

import (
	"hero_stroy/biz_server/base"
	"hero_stroy/biz_server/mod/user/userdata"
	"hero_stroy/comm/log"
)

const sqlGetUserByName = `SELECT user_id,user_name,password,hero_avatar FROM t_user WHERE user_name =?`

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
	if err := row.Scan(&user.UserId, &user.UserName, &user.Password, &user.HeroAvatar); err != nil {
		log.Error("GetUserByName error: %v", err)
		return nil
	}
	return user
}
