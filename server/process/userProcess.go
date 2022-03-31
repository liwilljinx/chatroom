package process

import (
	"chatroom/common/message"
	model2 "chatroom/common/model"
	"chatroom/common/util"
	"chatroom/server/model"
	"encoding/json"
	"fmt"
	"net"
)

// Process 处理用户相关的请求
type Process struct {
	Conn   net.Conn
	UserId int
}

// sendResMes 向客户端返回消息的方法
func (process Process) sendResMes(resMes interface{}, mesType string) (err error) {
	resMesJson, err := json.Marshal(resMes)
	if err != nil {
		fmt.Println("json.Marshal err=", err)
		return
	}

	mes := message.Mes{
		Type: mesType,
		Data: string(resMesJson),
	}

	tf := util.Transfer{
		Conn: process.Conn,
	}
	err = tf.WriteMes(&mes)
	return
}

// NotifyProcess 发送用户上线及离线的通知
func (process Process) NotifyProcess(userId int, userName string, userStatus int) {
	// 遍历在线用户列表，将新用户上线的消息通知给其他的在线用户
	for i, up := range UserMgrGlobal.UsersOnline {
		if i == userId {
			continue
		}
		var user = model2.User{
			UserId:     userId,
			UserName:   userName,
			UserStatus: userStatus,
		}
		var notifyOthersMes = message.NotifyOthersMes{
			User: user,
		}
		err := up.sendResMes(notifyOthersMes, message.NotifyOthersMesType)
		if err != nil {
			fmt.Println("NotifyProcess sendResMes err=", err)
			return
		}
	}
}

// RegisterProcess 处理注册的消息
func (process Process) RegisterProcess(data string) (err error) {
	var registerMes message.RegisterMes
	err = json.Unmarshal([]byte(data), &registerMes)
	if err != nil {
		fmt.Println("registerProcess json.Unmarshal err=", err)
		return
	}
	user := model2.User{}
	user = registerMes.User
	// 数据库验证是否注册成功
	err = model.MyUserDao.RegisterVerify(&user)

	resMes := message.ResMes{}
	if err != nil {
		// 用户存在，返回401错误
		if err == model.USER_EXIT {
			resMes.Code = 401
			resMes.Error = err.Error()
		} else {
			resMes.Code = 500
			resMes.Error = "服务器未知错误"
		}
	} else {
		// 将注册的新用户添加到AllUser中
		model.AllUser[user.UserId] = &user
		resMes.Code = 200
	}
	// 返回注册的结果
	err = process.sendResMes(resMes, message.ResMesType)
	return
}

// LoginProcess 处理登录
func (process Process) LoginProcess(data string) (err error) {
	var loginMes message.LoginMes
	err = json.Unmarshal([]byte(data), &loginMes)
	if err != nil {
		fmt.Println("json.Unmarshal err=", err)
		return
	}
	loginResMes := message.LoginResMes{}
	// 数据库验证是否登录成功
	user, err := model.MyUserDao.LoginVerify(loginMes.UserId, loginMes.UserPwd)
	if err != nil {
		if err == model.USER_NOT_EXIT { // 用户不存在
			loginResMes.ResMes.Code = 403
			loginResMes.ResMes.Error = err.Error()
		} else if err == model.PASSWORD_NOT_RIGHT { // 密码错误
			loginResMes.ResMes.Code = 400
			loginResMes.ResMes.Error = err.Error()
		} else {
			loginResMes.ResMes.Code = 500
			loginResMes.ResMes.Error = "服务器未知错误"
		}
	} else {
		fmt.Println(user, "登录了")
		process.UserId = loginMes.UserId
		loginMes.UserName = user.UserName
		// 登录成功后，更新AllUser中对应的用户状态
		model.ChangeUserStatus(loginMes.UserId, message.Online)
		// 新登录用户，添加到在线用户的map池里面
		UserMgrGlobal.AddUser(&process)
		// 将当前所有用户返回给客户端
		for _, v := range model.AllUser {
			onLineUser := &model2.User{
				UserId:     v.UserId,
				UserName:   v.UserName,
				UserStatus: v.UserStatus,
			}
			loginResMes.OnlineUsers = append(loginResMes.OnlineUsers, onLineUser)
		}
		// 向其他用户广播新用户登录
		process.NotifyProcess(loginMes.UserId, loginMes.UserName, message.Online)
		loginResMes.ResMes.Code = 200
	}
	// 发送响应的消息
	err = process.sendResMes(loginResMes, message.ResMesType)
	return
}

// SignOutProcess 处理下线的消息
func (process Process) SignOutProcess(data string) {
	var signOutMes message.SignOutMes
	err := json.Unmarshal([]byte(data), &signOutMes)
	if err != nil {
		fmt.Println("SignOutProcess json.Unmarshal err=", err)
		return
	}
	// 将该用户状态改为离线
	model.ChangeUserStatus(signOutMes.UserId, message.Offline)
	// 从在线用户中删除该用户
	UserMgrGlobal.DelUser(signOutMes.UserId)
	// 向其他用户发送该用户的下线消息
	process.NotifyProcess(signOutMes.UserId, signOutMes.UserName, message.Offline)
}
