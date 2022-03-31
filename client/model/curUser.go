package model

import (
	model2 "chatroom/common/model"
	"net"
)

// CurUser 保存当前登录的用户，以及连接
type CurUser struct {
	Conn net.Conn
	User model2.User
}
