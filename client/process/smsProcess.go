package process

import (
	"chatroom/common/message"
	"chatroom/common/util"
	"encoding/json"
	"fmt"
)

// SmsProcess 实现与聊天相关的功能
type SmsProcess struct {
}

// SendGroupMes 发送群聊的消息
func (smsProcess *SmsProcess) SendGroupMes(content string) {

	// 构建群聊的消息结构体实例
	smsMes := &message.SmsMes{
		Content: content,
		User:    MyCurUser.User,
		Type:    message.SmsGroupType,
	}

	smsJsonMes, err := json.Marshal(smsMes)
	if err != nil {
		fmt.Println("SendGroupMes json.Marshal err=", err)
		return
	}
	// 构建消息结构体实例
	mes := &message.Mes{
		Type: message.SmsMesType,
		Data: string(smsJsonMes),
	}
	// 创建一个读写消息的实例，这里是写操作
	tf := &util.Transfer{
		Conn: MyCurUser.Conn,
	}
	err = tf.WriteMes(mes)
	if err != nil {
		fmt.Println("SendGroupMes writeMes err=", err)
		return
	}
}

// SendUserMsg 发送私聊的消息
func (smsProcess *SmsProcess) SendUserMsg(userId int, content string) {
	// 构建私聊消息结构体实例
	smsMes := &message.SmsMes{
		Content:   content,
		User:      MyCurUser.User,
		RecUserId: userId,
		Type:      message.SmsPrivateType,
	}
	smsJsonMes, err := json.Marshal(smsMes)
	if err != nil {
		fmt.Println("SendGroupMes json.Marshal err=", err)
		return
	}
	// 构建消息结构体实例
	mes := &message.Mes{
		Type: message.SmsMesType,
		Data: string(smsJsonMes),
	}
	// 构建读写消息的实例，这里是写操作
	tf := &util.Transfer{
		Conn: MyCurUser.Conn,
	}
	err = tf.WriteMes(mes)
	if err != nil {
		fmt.Println("SendGroupMes writeMes err=", err)
		return
	}
}

// GetSmsList 登录时发送获取离线消息的请求
func (smsProcess SmsProcess) GetSmsList(userId int) {
	// 构建获取离线消息结构体的实例
	smsListMes := &message.SmsListMes{
		UserId: userId,
	}
	smsListJsonMes, err := json.Marshal(smsListMes)
	if err != nil {
		fmt.Println("GetSmsList json.Marshal err=", err)
		return
	}
	// 构建消息结构体的实例
	mes := &message.Mes{
		Type: message.SmsListMesType,
		Data: string(smsListJsonMes),
	}
	// 构建读写消息的实例，这里是写操作
	tf := &util.Transfer{
		Conn: MyCurUser.Conn,
	}
	err = tf.WriteMes(mes)
	if err != nil {
		fmt.Println("GetSmsList writeMes err=", err)
		return
	}
}
