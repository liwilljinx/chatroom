package main

import (
	"chatroom/common/message"
	"chatroom/common/util"
	"chatroom/server/process"
	"fmt"
	"io"
	"net"
)

// TransProcess 用户处理各种类型的客户端请求
type TransProcess struct {
	Conn net.Conn
}

// ProcessMes 根据不同类型的客户端请求，调用不同的处理函数
func (transProcess TransProcess) ProcessMes(mes *message.Mes) (err error) {
	data := mes.Data
	switch mes.Type {
	case message.LoginMesType:
		lp := &process.Process{
			Conn: transProcess.Conn,
		}
		// 登录
		err = lp.LoginProcess(data)
	case message.RegisterMesType:
		lp := &process.Process{
			Conn: transProcess.Conn,
		}
		// 注册
		err = lp.RegisterProcess(data)
	case message.SmsMesType:
		sp := &process.SmsProcess{}
		// 消息转发
		sp.TransmitMes(data)
	case message.SmsListMesType:
		sp := &process.SmsProcess{}
		// 获取离线消息
		sp.SendListSms(data)
	case message.SignOutMesType:
		lp := &process.Process{
			Conn: transProcess.Conn,
		}
		// 退出登录
		lp.SignOutProcess(data)
	default:
		fmt.Println("消息类型错误，没有这种消息类型")
	}
	return
}

// TransMes 等待客户端发送请求
func (transProcess TransProcess) TransMes() {
	tf := &util.Transfer{
		Conn: transProcess.Conn,
	}
	for {
		mes, err := tf.ReadMes()
		if err != nil {
			if err == io.EOF {
				fmt.Println("客户端退出")
				return
			}
			if _, ok := err.(*net.OpError); ok {
				fmt.Println("客户端退出")
				return
			}
			fmt.Println("接受数据出错，err=", err)
			return
		}
		err = transProcess.ProcessMes(&mes)
		if err != nil {
			fmt.Println("处理消息出错，err=", err)
			return
		}
	}
}
