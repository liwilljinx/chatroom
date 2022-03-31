package message

import "chatroom/common/model"

// 不同消息的类型，用于客户端和服务器之间的消息通讯
const (
	LoginMesType        = "LoginMes"
	RegisterMesType     = "RegisterMes"
	ResMesType          = "ResMes"
	NotifyOthersMesType = "NotifyOthersMes"
	SmsMesType          = "SmsMes"
	SmsResMesType       = "SmsResMes"
	SignOutMesType      = "SignOutMes"
	SmsListMesType      = "SmsListMes"
	SmsListResMesType   = "SmsListResMes"
)

// 用户当前的状态，Offline为0，Online为1
const (
	Offline = iota
	Online
)

// 发送消息的类型，群聊或者私聊
const (
	SmsGroupType   = "Group"
	SmsPrivateType = "Private"
)

// Mes 发送消息的结构体，下面各种类型的消息封装好后，保存到Data里面
type Mes struct {
	Type string `json:"type"`
	Data string `json:"data"`
}

// LoginMes 客户端用于登录的消息结构体
type LoginMes struct {
	UserId   int    `json:"userId"`
	UserPwd  string `json:"userPwd"`
	UserName string `json:"userName"`
}

// RegisterMes 客户端用于注册的消息结构体
type RegisterMes struct {
	User model.User `json:"user"`
}

// ResMes 服务器响应客户端注册和登录消息的结构体
type ResMes struct {
	Code  int    `json:"code"`
	Error string `json:"error"`
}

// LoginResMes 服务器响应登录消息的结构体，OnlineUsers记录服务器中所有用户的信息，将这个返回给客户端
type LoginResMes struct {
	ResMes      ResMes        `json:"resMes"`
	OnlineUsers []*model.User `json:"onlineUsers"`
}

// SignOutMes 客户端用于发送下线消息的结构体
type SignOutMes struct {
	UserId   int    `json:"userId"`
	UserName string `json:"userName"`
}

// NotifyOthersMes 服务器用于发送用户状态消息的结构体
type NotifyOthersMes struct {
	User model.User `json:"user"`
}

// SmsMes 客户端用于发送消息的结构体
type SmsMes struct {
	Content   string     `json:"content"`
	User      model.User `json:"user"`
	Type      string     `json:"type"`
	RecUserId int        `json:"recUserId"`
}

// SmsResMes 服务器用于将客户端发送的消息转发至其他用户的结构体
type SmsResMes struct {
	Content string     `json:"content"`
	User    model.User `json:"user"`
	Type    string     `json:"type"`
}

// SmsListMes 客户端获取离线消息的结构体
type SmsListMes struct {
	UserId int `json:"userId"`
}

// SmsListResMes 服务器返回离线消息的结构体
type SmsListResMes struct {
	SmsList []string `json:"smsList"`
}
