package model

import model2 "chatroom/common/model"

// AllUser 保存所有的用户信息
var (
	AllUser map[int]*model2.User
)

// 初始化AllUser
func init() {
	AllUser = make(map[int]*model2.User, 1024)
}

// ChangeUserStatus 改变用户的状态，在线或离线
func ChangeUserStatus(userId int, userStatus int) {
	user, _ := AllUser[userId]
	user.UserStatus = userStatus
}
