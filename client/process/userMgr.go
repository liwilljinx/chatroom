package process

import (
	"chatroom/client/model"
	"chatroom/common/message"
	model2 "chatroom/common/model"
	"fmt"
)

// AllUsers 用户保存所有用户
var AllUsers = make(map[int]*model2.User, 10)

// MyCurUser 记录当前登录的用户，包括用户的信息和客户端与服务器的连接
var MyCurUser model.CurUser

// OutputOnlineUsers 显示当前在线的用户
func OutputOnlineUsers() {
	flag := 0
	for _, user := range AllUsers {
		if flag == 0 {
			fmt.Println("当前在线用户")
		}
		switch user.UserStatus {
		case message.Online:
			fmt.Printf("%s[在线]\n", user.UserName)
		case message.Offline:
			fmt.Printf("%s[离线]\n", user.UserName)
		}
		flag++
	}
}

// UpdateUserStatus 更新用户的状态，在线或者是离线；如果是新注册的用户那么会添加一个新的用户到AllUsers中
func UpdateUserStatus(notifyOthersMes *message.NotifyOthersMes) {
	user := &notifyOthersMes.User
	AllUsers[user.UserId] = user
	if user.UserStatus == message.Online {
		fmt.Printf("\n%s上线了\n", user.UserName)
	} else {
		fmt.Printf("\n%s下线了\n", user.UserName)
	}
}
