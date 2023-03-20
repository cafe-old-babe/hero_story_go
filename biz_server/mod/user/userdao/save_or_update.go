package userdao

import (
	"hero_story/biz_server/base"
	"hero_story/biz_server/mod/user/userdata"
	"hero_story/comm/log"
)

const sqlSaveOrUpdate = `
insert  into t_user(user_name, password, hero_avatar, create_time, last_login_time) values (?, ?, ?, ?, ?)
on duplicate key update last_login_time=?
`

func SaveOrUpdate(user *userdata.User) error {
	if user == nil {
		return nil
	}
	stmt, err := base.MysqlDB.Prepare(sqlSaveOrUpdate)
	if err != nil {
		log.Error("saveOrUpdate Prepare err: %+v", err)
		return err
	}
	exec, err := stmt.Exec(user.UserName, user.Password, user.HeroAvatar, user.CreateTime, user.LastLoginTime, user.LastLoginTime)
	if err != nil {
		log.Error("saveOrUpdate Exec err: %+v", err)
		return err
	}
	id, err := exec.LastInsertId()
	if err != nil {
		log.Error("saveOrUpdate LastInsertId err: %+v", err)
		return err
	}
	user.UserId = id
	return nil
}
