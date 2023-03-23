package userdata

import "sync"

type User struct {
	UserId        int64  `db:"user_id"`
	UserName      string `db:"user_name"`
	Password      string `db:"password"`
	HeroAvatar    string `db:"hero_avatar"`
	CurrHp        int32  `db:"curr_hp"`
	CreateTime    int64  `db:"create_time"`
	LastLoginTime int64  `db:"last_login_time"`
	MoveState     *MoveState

	componentMap  *sync.Map
	createMapLock sync.Mutex
}

func (u *User) GetComponentMap() *sync.Map {

	if u.componentMap != nil {
		goto mapLabel

	}
	u.createMapLock.Lock()
	defer u.createMapLock.Unlock()

	if u.componentMap != nil {
		u.createMapLock.Unlock()
		goto mapLabel
	}
	u.componentMap = &sync.Map{}
mapLabel:
	return u.componentMap

}
