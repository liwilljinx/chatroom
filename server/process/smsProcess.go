package process

import (
	"chatroom/common/message"
	"chatroom/common/util"
	"chatroom/server/model"
	"encoding/json"
	"fmt"
	"net"
)

// SmsProcess 处理消息相关的请求
type SmsProcess struct {
}

// TransmitMes 转发消息
func (smsProcess *SmsProcess) TransmitMes(data string) {
	smsMes := message.SmsMes{}
	err := json.Unmarshal([]byte(data), &smsMes)
	if err != nil {
		fmt.Println("TransmitMes json.Unmarshal err=", err)
		return
	}
	// 创建响应消息的结构体
	smsResMes := &message.SmsResMes{
		User:    smsMes.User,
		Content: smsMes.Content,
		Type:    smsMes.Type,
	}
	smsResJsonMes, err := json.Marshal(smsResMes)
	if err != nil {
		fmt.Println("TransmitMes json.Marshal err=", err)
		return
	}
	mes := message.Mes{
		Type: message.SmsResMesType,
		Data: string(smsResJsonMes),
	}
	// 根据不同的消息类型进行转发
	if smsMes.Type == message.SmsGroupType {
		// 群聊，转发给所有人
		smsProcess.SendToGroupMes(smsMes.User.UserId, &mes)
	} else {
		// 私聊，转发给指定的用户
		up, ok := UserMgrGlobal.UsersOnline[smsMes.RecUserId]
		if !ok {
			fmt.Println("该用户不在线")
			// 如果用户不在线，将消息保存到数据库中
			smsProcess.SaveOffLineSms(smsMes.RecUserId, mes.Data)
			return
		}
		smsProcess.SendToOtherMes(up.Conn, &mes)
	}
}

// SaveOffLineSms 如果用户不在线，那就把消息保存到redis中
func (smsProcess *SmsProcess) SaveOffLineSms(userId int, data string) {
	err := model.MySmsDao.SaveSms(userId, data)
	if err != nil {
		fmt.Println("SaveOffLineSms SaveSms err=", err)
		return
	}
}

// SendToGroupMes 群发消息
func (smsProcess *SmsProcess) SendToGroupMes(userId int, mes *message.Mes) {
	for _, v := range model.AllUser {
		// 过滤掉自己
		if v.UserId == userId {
			continue
		}
		// 根据用户状态，判断是否在线，若在线，直接发送，不在线，保存到数据库中
		if v.UserStatus == message.Online {
			up, _ := UserMgrGlobal.UsersOnline[v.UserId]
			smsProcess.SendToOtherMes(up.Conn, mes)
		} else {
			smsProcess.SaveOffLineSms(v.UserId, mes.Data)
		}
	}
}

// SendToOtherMes 发送消息的方法
func (smsProcess *SmsProcess) SendToOtherMes(conn net.Conn, mes *message.Mes) {
	tf := &util.Transfer{
		Conn: conn,
	}
	err := tf.WriteMes(mes)
	if err != nil {
		fmt.Println("SendToOtherMes write err=", err)
		return
	}
}

// SendListSms 发送离线消息
func (smsProcess SmsProcess) SendListSms(data string) {
	smsListMes := message.SmsListMes{}
	err := json.Unmarshal([]byte(data), &smsListMes)
	if err != nil {
		fmt.Println("SendListSms json.Unmarshal err=", err)
		return
	}
	smsList := model.MySmsDao.GetSms(smsListMes.UserId)
	smsListResMes := message.SmsListResMes{
		SmsList: smsList,
	}
	smsListResMesJson, err := json.Marshal(smsListResMes)
	if err != nil {
		fmt.Println("SendListSms json.Marshal err=", err)
		return
	}
	mes := message.Mes{
		Type: message.SmsListResMesType,
		Data: string(smsListResMesJson),
	}
	// 从在线用户列表中取出该客户端的连接，并将消息发送过去
	up, _ := UserMgrGlobal.UsersOnline[smsListMes.UserId]
	smsProcess.SendToOtherMes(up.Conn, &mes)
}
