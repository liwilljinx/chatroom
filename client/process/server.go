package process

import (
	"bufio"
	"chatroom/common/message"
	"chatroom/common/util"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"strings"
)

// ShowMenu 登录成功后展示的页面
func ShowMenu(userName string) (key int) {
	fmt.Printf("\n--------------恭喜%s登录成功-------------\n", userName)
	fmt.Println("--------------1.显示在线用户列表--------------")
	fmt.Println("--------------2.发送消息--------------------")
	fmt.Println("--------------3.信息列表--------------------")
	fmt.Println("--------------4.退出系统--------------------")
	fmt.Print("--------------请输入（1-4）:")
	fmt.Scanln(&key)
	switch key {
	case 1:
		fmt.Println("显示用户列表")
		OutputOnlineUsers()
	case 2:
		var smsKey int
		fmt.Print("发送消息[群聊->1]/[私聊->2]:")
		fmt.Scanln(&smsKey)
		ShowSms(smsKey)
	case 3:
		fmt.Println("信息列表")
		ShowSmsList()
	case 4:
		fmt.Println("退出")
	default:
		fmt.Println("输入不正确")
	}
	return
}

// ShowSms 处理群聊和私聊两种消息的输入
func ShowSms(smsKey int) {
	var content string
	smsProcess := &SmsProcess{}
	if smsKey == 1 {
		fmt.Print("请输入你想对大家说的话:")
		reader := bufio.NewReader(os.Stdin)  // 标准输入输出
		content, _ = reader.ReadString('\n') // 回车结束
		content = strings.TrimSpace(content)
		// 调用群发消息的功能
		smsProcess.SendGroupMes(content)
	} else {
		var ToUserId int
		fmt.Print("请输入发送的用户id:")
		fmt.Scanln(&ToUserId)
		fmt.Print("请输入你想对TA说的话:")
		reader := bufio.NewReader(os.Stdin)  // 标准输入输出
		content, _ = reader.ReadString('\n') // 回车结束
		content = strings.TrimSpace(content)
		// 调用私聊的功能
		smsProcess.SendUserMsg(ToUserId, content)
	}
}

// processServer 登录成功后保持与服务器的连接，处理服务器发送过来不同类型的消息种类
func processServer(conn net.Conn) {
	tf := util.Transfer{
		Conn: conn,
	}
	for {
		// 等待服务器发送消息，如果服务器还没发送消息回来，这里会一直阻塞，直到服务器发送消息或者连接由某一方断开而产生错误
		mes, err := tf.ReadMes()
		if err != nil {
			if _, ok := err.(*net.OpError); ok {
				fmt.Println("退出系统")
				return
			}
			fmt.Println("readMes err=", err)
			return
		}
		switch mes.Type {
		// 用户上线的消息类型
		case message.NotifyOthersMesType:
			var notifyOthersMes = message.NotifyOthersMes{}
			err = json.Unmarshal([]byte(mes.Data), &notifyOthersMes)
			if err != nil {
				fmt.Println("message.NotifyOthersMes json.Unmarshal err=", err)
				break
			}
			// 当接收到用户登录或者下线的消息时，就更新AllUsers中该用户的状态
			UpdateUserStatus(&notifyOthersMes)
		// 接收聊天消息的消息类型
		case message.SmsResMesType:
			var smsResMes = message.SmsResMes{}
			err = json.Unmarshal([]byte(mes.Data), &smsResMes)
			if err != nil {
				fmt.Println("message.SmsResMes json.Unmarshal err=", err)
				break
			}
			// 将接收到的消息保存到smsList中
			SaveSmsList(&smsResMes)
			// 打印出该消息的内容
			OutPutSms(&smsResMes)
		// 接收离线消息的消息类型
		case message.SmsListResMesType:
			var smsListResMes = message.SmsListResMes{}
			err = json.Unmarshal([]byte(mes.Data), &smsListResMes)
			if err != nil {
				fmt.Println("message.SmsListResMesType json.Unmarshal err=", err)
				break
			}
			// 将离线的时收到的留言保存到smsList中
			SaveSmsList(&smsListResMes)
		}
	}
}
