package process

import (
	"chatroom/server/model"
)

// UserMgrGlobal 创建管理用户的全局变量
var (
	UserMgrGlobal *UserMgr
)

// UserMgr 用于管理在线用户，保存用户的id和客户端连接
type UserMgr struct {
	UsersOnline map[int]*Process
}

//初始化
func init() {
	UserMgrGlobal = &UserMgr{
		UsersOnline: make(map[int]*Process, 1024),
	}
}

// AddUser 增加在线用户
func (userMgr *UserMgr) AddUser(up *Process) {
	userMgr.UsersOnline[up.UserId] = up
}

// DelUser 删除在线用户
func (userMgr *UserMgr) DelUser(userId int) {
	delete(userMgr.UsersOnline, userId)
}

// GetAllOnlineUser 返回在线用户
func (userMgr *UserMgr) GetAllOnlineUser() map[int]*Process {
	return userMgr.UsersOnline
}

// GetUserById 用过id获取在线用户
func (userMgr *UserMgr) GetUserById(userId int) (up *Process, err error) {
	up, ok := userMgr.UsersOnline[userId]
	if !ok {
		err = model.USER_NOT_ONLINE
		return
	}
	return
}
