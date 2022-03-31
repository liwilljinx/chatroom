package model

import (
	"fmt"
	"github.com/garyburd/redigo/redis"
)

// MySmsDao 全局的smsDao，用于管理与redis连接，处理聊天类型的数据请求
var MySmsDao *SmsDao

// SmsDao 用于管理离线消息的结构体
type SmsDao struct {
	Pool *redis.Pool
}

// NewSmsDao 初始化一个Sms连接池
func NewSmsDao(pool *redis.Pool) (smsDao *SmsDao) {
	smsDao = &SmsDao{
		Pool: pool,
	}
	return
}

// SaveSms 当用户不在线的时候，保存消息
func (smsDao SmsDao) SaveSms(userId int, data string) (err error) {
	conn := smsDao.Pool.Get()
	defer conn.Close()

	_, err = conn.Do("LPUSH", userId, data)
	if err != nil {
		fmt.Println("SaveSms conn.Do err=", err)
		return
	}
	return
}

// GetSms 用户上线后从redis中取出消息
func (smsDao SmsDao) GetSms(userId int) []string {
	conn := smsDao.Pool.Get()
	defer conn.Close()

	// 将该用户的留言取出，redis中不再保存聊天信息
	var smsDatas []string
	for {
		res, err := redis.String(conn.Do("RPOP", userId))
		if err != nil {
			if err == redis.ErrNil {
				break
			}
			fmt.Println("GetSms conn.Do err=", err)
			continue
		}
		smsDatas = append(smsDatas, res)
	}
	return smsDatas
}
