package model

// User 保存用户信息的结构体
type User struct {
	UserId     int    `json:"userId"`     // 用户id，唯一
	UserPwd    string `json:"userPwd"`    // 用户密码
	UserName   string `json:"userName"`   // 用户名
	UserStatus int    `json:"userStatus"` // 用户状态，在线或离线
}
