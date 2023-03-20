package userdata

type userGroup struct {
	innerMap map[int64]*User
}

var userGroupInstance = &userGroup{
	innerMap: make(map[int64]*User),
}

// GetUserGroup 获取用户组信息
func GetUserGroup() *userGroup {

	return userGroupInstance
}

// Add 添加用户
func (u *userGroup) Add(user *User) {
	if nil == user {
		return
	}
	u.innerMap[user.UserId] = user
}

// RemoveByUserId 删除用户
func (u *userGroup) RemoveByUserId(userId int64) {
	if userId <= 0 {
		return
	}
	delete(u.innerMap, userId)
}

// GetByUserId 根据用户ID获取用户
func (group *userGroup) GetByUserId(userId int64) (user *User) {
	if userId <= 0 {
		return
	}
	user = group.innerMap[userId]
	return
}
