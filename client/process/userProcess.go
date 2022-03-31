package process

import (
	"chatroom/common/message"
	"chatroom/common/model"
	"chatroom/common/util"
	"encoding/json"
	"fmt"
	"net"
)

// Process 实现与用户相关的功能；注册、登录、下线
type Process struct {
}

// Register 用于注册用户
func (process Process) Register(userId int, userPwd, userName string) {
	// 创建与服务器的连接
	conn, err := net.Dial("tcp", "localhost:8888")
	defer conn.Close()

	if err != nil {
		fmt.Println("register net.Dial err=", err)
		return
	}
	// 创建用户结构体实例
	user := model.User{
		UserId:   userId,
		UserPwd:  userPwd,
		UserName: userName,
	}
	registerMes := message.RegisterMes{
		User: user,
	}
	registerJsonMes, _ := json.Marshal(registerMes)
	// 创建消息结构体实例
	mes := message.Mes{
		Type: message.RegisterMesType,
		Data: string(registerJsonMes),
	}

	tf := &util.Transfer{
		Conn: conn,
	}
	// 发送注册消息
	err = tf.WriteMes(&mes)
	if err != nil {
		fmt.Println("register tf.WriterMes err=", err)
		return
	}

	// 等待服务器返回的注册响应
	res, err := tf.ReadMes()
	resData := message.ResMes{}
	err = json.Unmarshal([]byte(res.Data), &resData)
	if err != nil {
		fmt.Println("register json.Unmarshal err=", err)
		return
	}
	if resData.Code == 200 {
		fmt.Println("注册成功，请登录")
	} else {
		fmt.Println(resData.Error)
	}
}

// Login 用户登录
func (process Process) Login(userId int, userPwd string) {
	// 创建连接
	conn, err := net.Dial("tcp", "localhost:8888")
	defer conn.Close()
	if err != nil {
		fmt.Println("net.Dial err=", err)
		return
	}
	// 创建登录结构体实例
	loginMes := message.LoginMes{
		UserId:  userId,
		UserPwd: userPwd,
	}
	loginMesJsonData, _ := json.Marshal(loginMes)
	// 创建消息结构体实例
	mes := message.Mes{
		Type: message.LoginMesType,
		Data: string(loginMesJsonData),
	}

	tf := &util.Transfer{
		Conn: conn,
	}
	// 发送登录的消息
	err = tf.WriteMes(&mes)

	if err != nil {
		fmt.Println("tf.Writes err=", err)
		return
	}
	// 等待服务器返回的登录响应
	resMes, err := tf.ReadMes()
	loginResMes := message.LoginResMes{}
	err = json.Unmarshal([]byte(resMes.Data), &loginResMes)
	if err != nil {
		fmt.Println("json.Unmarshal, err=", err)
		return
	}
	// 如果状态码是200，那么登录成功
	if loginResMes.ResMes.Code == 200 {

		// 将当前的用户id、用户名、用户状态、与服务器的连接保存到MyCurUser中，以方便其他功能的使用
		MyCurUser.Conn = conn
		MyCurUser.User.UserId = userId
		MyCurUser.User.UserStatus = message.Online
		// 启动一个监听服务器请求的协程，处理各种服务器返回的消息
		go processServer(conn)

		// 登录成功后，服务器会返回所有用户的信息，将用户信息保存到AllUsers中，用于对所有用户的管理
		for _, v := range loginResMes.OnlineUsers {
			if v.UserId == userId {
				MyCurUser.User.UserName = v.UserName
				continue
			}
			switch v.UserStatus {
			case message.Online:
				fmt.Printf("%s[在线]\n", v.UserName)
			case message.Offline:
				fmt.Printf("%s[离线]\n", v.UserName)
			}
			user := &model.User{
				UserId:     v.UserId,
				UserName:   v.UserName,
				UserStatus: v.UserStatus,
			}
			AllUsers[v.UserId] = user
		}
		// 登录成功后，向服务器发送同步离线消息的请求
		smsProcess := &SmsProcess{}
		smsProcess.GetSmsList(userId)
		for {
			// 展示登录成功后的功能页面
			key := ShowMenu(MyCurUser.User.UserName)
			if key == 4 {
				// 向服务器发送下线的消息
				process.SignOut(userId, MyCurUser.User.UserName, conn)
				break
			}
		}

	} else {
		fmt.Println("err=", loginResMes.ResMes.Error)
	}
	return
}

// SignOut 退出系统，向服务器发送退出的消息
func (process Process) SignOut(userId int, userName string, conn net.Conn) {
	// 创建离线消息的结构体
	signOutMes := message.SignOutMes{
		UserId:   userId,
		UserName: userName,
	}
	signOutJsonMes, err := json.Marshal(signOutMes)
	if err != nil {
		fmt.Println("SignOut json.Marshal err=", err)
		return
	}
	mes := message.Mes{
		Type: message.SignOutMesType,
		Data: string(signOutJsonMes),
	}

	tf := &util.Transfer{
		Conn: conn,
	}
	// 发送该消息
	err = tf.WriteMes(&mes)
	if err != nil {
		fmt.Println("SignOut tf.WriterMes err=", err)
		return
	}
}
