package model

import "errors"

// 自定义错误
var (
	USER_NOT_EXIT      = errors.New("用户不存在")
	PASSWORD_NOT_RIGHT = errors.New("账户或者密码错误")
	USER_EXIT          = errors.New("用户已经存在")
	USER_NOT_ONLINE    = errors.New("用户不在线")
)
