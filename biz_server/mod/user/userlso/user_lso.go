package userlso

import (
	"fmt"
	"hero_story/biz_server/mod/user/userdao"
	"hero_story/biz_server/mod/user/userdata"
	"hero_story/comm/async_op"
	"sync"
)

type myFunc = func()
type UserLso struct {
	*userdata.User
}

var asyncOp myFunc
var locker = &sync.Mutex{}

func (user *UserLso) GetLsoId() string {
	return fmt.Sprintf("UserLso_%d", user.UserId)
}

func (user *UserLso) SaveOrUpdate() {
	if asyncOp != nil {
		goto doLabel
	}
	locker.Lock()
	if asyncOp == nil {
		asyncOp = func() {
			userdao.SaveOrUpdate(user.User)
		}
		locker.Unlock()

	}
doLabel:
	async_op.Process(int(user.UserId), &asyncOp, nil)

}
