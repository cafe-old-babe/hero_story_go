package user_dao

import (
	"hero_story/biz_server/base"
	"hero_story/biz_server/mod/user/user_data"
	"hero_story/comm/log"
)

const sqlSaveOrUpdate = `
insert  into t_user(user_name, password, hero_avatar, curr_hp, create_time, last_login_time) values (?, ?, ?, ?, ?, ?)
on duplicate key update curr_hp = values(curr_hp), last_login_time = values(last_login_time)
`

func SaveOrUpdate(user *user_data.User) error {
	if user == nil {
		return nil
	}
	stmt, err := base.MysqlDB.Prepare(sqlSaveOrUpdate)
	if err != nil {
		log.Error("saveOrUpdate Prepare err: %+v", err)
		return err
	}
	exec, err := stmt.Exec(user.UserName, user.Password, user.HeroAvatar, user.CurrHp, user.CreateTime, user.LastLoginTime)
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
