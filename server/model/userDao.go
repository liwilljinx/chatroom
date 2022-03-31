package model

import (
	"chatroom/common/model"
	"encoding/json"
	"fmt"
	"github.com/garyburd/redigo/redis"
)

// MyUserDao 全局的UserDao,用于处理用户相关的数据请求
var MyUserDao *UserDao

type UserDao struct {
	Pool *redis.Pool
}

// NewUserDao 初始化
func NewUserDao(pool *redis.Pool) (userDao *UserDao) {
	userDao = &UserDao{
		Pool: pool,
	}
	return
}

// LoadAllUser 加载所有用户的信息
func (userDao UserDao) LoadAllUser() {
	conn := userDao.Pool.Get()
	defer conn.Close()

	res, err := redis.Values(conn.Do("HVals", "users"))
	if err != nil {
		fmt.Println("LoadAllUser conn.Do err=", err)
		return
	}

	for _, v := range res {
		var user model.User
		err = json.Unmarshal(v.([]byte), &user)
		if err != nil {
			fmt.Printf("LoadAllUser json.Unmarshal %s err=%s\n", string(v.([]byte)), err)
			continue
		}
		AllUser[user.UserId] = &user
	}
}

// GetUserById 根据用户的id获取用户信息
func (userDao UserDao) GetUserById(conn redis.Conn, id int) (user model.User, err error) {
	res, err := redis.String(conn.Do("HGet", "users", id))
	if err != nil {
		// 如果没有找到该用户，则用户还没有注册，会抛出USER_NOT_EXIT异常
		if err == redis.ErrNil {
			err = USER_NOT_EXIT
			return
		}
		return
	}
	err = json.Unmarshal([]byte(res), &user)
	if err != nil {
		fmt.Println("userDao json.Unmarshal err=", err)
		return
	}
	return
}

// LoginVerify 验证用户的登录信息
func (userDao UserDao) LoginVerify(userId int, userPwd string) (user model.User, err error) {
	conn := userDao.Pool.Get()
	defer conn.Close()

	// 通过id查找用户，出错则表示用户不存在
	user, err = userDao.GetUserById(conn, userId)
	if err != nil {
		return
	}
	// 若存在用户，则比对密码，密码一致，登录成功
	if userPwd != user.UserPwd {
		err = PASSWORD_NOT_RIGHT
		return
	}
	return
}

// RegisterVerify 注册用户
func (userDao UserDao) RegisterVerify(user *model.User) (err error) {
	conn := userDao.Pool.Get()
	defer conn.Close()

	// 根据id获取用户，若用户存在，即没有发生错误，那么注册失败
	_, err = userDao.GetUserById(conn, user.UserId)
	if err == nil {
		err = USER_EXIT
		return
	}

	userJson, _ := json.Marshal(user)
	_, err = conn.Do("HSet", "users", user.UserId, userJson)
	if err != nil {
		fmt.Println("register conn.Do err=", err)
		return
	}
	return
}
