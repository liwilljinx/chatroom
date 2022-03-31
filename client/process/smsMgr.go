package process

import (
	"chatroom/common/message"
	"encoding/json"
	"fmt"
)

// smsList 用于保存发送过来的消息
var smsList []*message.SmsResMes

// OutPutSms 根据消息的不同类型打印出消息内容
func OutPutSms(smsResMes *message.SmsResMes) {
	if smsResMes.Type == message.SmsGroupType {
		fmt.Printf("\n%s对所有人说：%s\n", smsResMes.User.UserName, smsResMes.Content)
	} else {
		fmt.Printf("\n%s对你说：%s\n", smsResMes.User.UserName, smsResMes.Content)
	}
}

// SaveSmsList 将发送过来的消息保存到smsList中；有两种消息，在线时收到的消息和同步过来的离线消息
func SaveSmsList(data interface{}) {
	switch t := data.(type) {
	case *message.SmsResMes:
		smsList = append(smsList, t)
	case *message.SmsListResMes:
		for _, v := range t.SmsList {
			// 因为SmsListResMes的SmsList中的值是字符串类型，所以需要反序列化为对应的结构体
			smsResMes := message.SmsResMes{}
			err := json.Unmarshal([]byte(v), &smsResMes)
			if err != nil {
				fmt.Println("SaveSmsList json.Unmarshal err=", err)
				continue
			}
			smsList = append(smsList, &smsResMes)
		}
	default:
		fmt.Println("Unexpect type")
	}
}

// ShowSmsList 显示smsList中的消息
func ShowSmsList() {
	for _, v := range smsList {
		OutPutSms(v)
	}
}
