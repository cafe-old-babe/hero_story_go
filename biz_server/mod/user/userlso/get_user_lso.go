package userlso

import "hero_story/biz_server/mod/user/userdata"

// GetUserLso 避免频繁创建对象,产生内存碎片
func GetUserLso(user *userdata.User) *UserLso {
	if nil == user {
		return nil
	}
	existMap, _ := user.GetComponentMap().Load("UserLso")
	if nil != existMap {
		return existMap.(*UserLso)
	}

	existMap = &UserLso{
		User: user,
	}
	existMap, _ = user.GetComponentMap().LoadOrStore("UserLso", existMap)
	return existMap.(*UserLso)
}
