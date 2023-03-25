package user_lso

import (
	"fmt"
	"hero_story/biz_server/mod/user/user_dao"
	"hero_story/biz_server/mod/user/user_data"
	"hero_story/comm/async_op"
)

type UserLso struct {
	*user_data.User
}

func (user *UserLso) GetLsoId() string {
	return fmt.Sprintf("UserLso_%d", user.UserId)
}

func (user *UserLso) SaveOrUpdate(callback func()) {

	asyncOp := func() {
		user_dao.SaveOrUpdate(user.User)
		if callback != nil {
			callback()
		}
	}

	async_op.Process(int(user.UserId), &asyncOp, nil)

}
