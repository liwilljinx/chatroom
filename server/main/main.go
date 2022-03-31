package main

import (
	"chatroom/server/model"
	"fmt"
	"github.com/garyburd/redigo/redis"
	"net"
	"time"
)

// firstProcess 监听连接的函数
func firstProcess(conn net.Conn) {
	defer conn.Close()
	fmt.Println("等待输入")
	pms := TransProcess{
		Conn: conn,
	}
	pms.TransMes()
}

// redis连接池
var (
	pool *redis.Pool
)

// InitPool 初始化redis连接池
func InitPool(maxIdle int, maxActive int, idleTimeout time.Duration, address string) {
	pool = &redis.Pool{
		MaxIdle:     maxIdle,
		MaxActive:   maxActive,
		IdleTimeout: idleTimeout,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", address)
			if err != nil {
				return nil, err
			}
			c.Do("SELECT", 3)
			return c, nil
		},
	}
}

// InitDao 初始化两个dao层，UserDao和SmsDao
func InitDao(pool *redis.Pool) {
	model.MyUserDao = model.NewUserDao(pool)
	model.MySmsDao = model.NewSmsDao(pool)
}

func main() {

	InitPool(32, 0, time.Second*100, "localhost:6379")
	InitDao(pool)
	// 服务器启动的时候，加载redis中所有用户的信息
	model.MyUserDao.LoadAllUser()

	// 监听8888端口
	listen, err := net.Listen("tcp", "localhost:8888")
	defer listen.Close()
	if err != nil {
		fmt.Println("net.Listen err=", err)
		return
	}
	fmt.Println("正在监听8888端口")
	for {
		// 等待客户端连接
		conn, err := listen.Accept()
		if err != nil {
			fmt.Println("listen.Accept err=", err)
			continue
		}
		go firstProcess(conn)
	}
}
