package main

import (
	"chatroom/client/process"
	"fmt"
)

var userId int
var userPwd string
var userName string

func main() {

	var key int

	var loop = true

	for loop {
		fmt.Println("\n-------------------欢迎使用多人聊天室-------------------")
		fmt.Println("\t\t\t 1.登录系统")
		fmt.Println("\t\t\t 2.注册系统")
		fmt.Println("\t\t\t 3.退出系统")
		fmt.Printf("\t\t\t (请选择1-3):")
		fmt.Scanln(&key)
		switch key {
		case 1:
			fmt.Println("登录")
			fmt.Print("请输入用户Id:")
			fmt.Scanln(&userId)
			fmt.Print("请输入密码:")
			fmt.Scanln(&userPwd)
			ps := process.Process{}
			ps.Login(userId, userPwd)
			loop = false
		case 2:
			fmt.Println("注册")
			fmt.Print("请输入用户Id:")
			fmt.Scanln(&userId)
			fmt.Print("请输入密码:")
			fmt.Scanln(&userPwd)
			fmt.Print("请输入用户名:")
			fmt.Scanln(&userName)
			ps := process.Process{}
			ps.Register(userId, userPwd, userName)
		case 3:
			loop = false
		default:
			fmt.Println("输入有误，请重新输入")
		}
	}
}
