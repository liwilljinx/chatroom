package util

import (
	"chatroom/common/message"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
)

// Transfer 客户端和服务器进行通讯时的工具包
type Transfer struct {
	Conn net.Conn       // 客户端与服务器的连接
	Buf  [1024 * 4]byte // 用于消息的缓存
}

func (transfer *Transfer) ReadMes() (mes message.Mes, err error) {
	// 读取消息体的长度大小
	_, err = transfer.Conn.Read(transfer.Buf[:4])
	if err != nil {
		if err == io.EOF {
			return
		}
		return
	}

	var pkgLen uint32
	// 转换消息体的长度大小
	pkgLen = binary.BigEndian.Uint32(transfer.Buf[:4])
	// 读取发送过来的消息体，返回的n就是消息体的长度
	// 因为这里会阻塞读取消息，所以对应的ReadMes方法与WriteMes方法需要有成对的Read和Write操作
	n, err := transfer.Conn.Read(transfer.Buf[:pkgLen])

	// 如果第一次接收的消息体长度大小和第二次接收的消息体的长度n一致，那么传输过程中没有丢包，如果不一致，则产生了丢包
	if n != int(pkgLen) || err != nil {
		err = errors.New("read body err")
		return
	}
	// 将消息反序列化为消息类型的结构体
	err = json.Unmarshal(transfer.Buf[:pkgLen], &mes)
	if err != nil {
		err = errors.New("unmarshal pkg err")
		return
	}
	return
}

func (transfer *Transfer) WriteMes(mes *message.Mes) (err error) {
	// 将消息类型的结构体序列化
	mesJsonData, _ := json.Marshal(mes)

	// 计算消息的长度大小，并发送
	var pkgLen uint32
	pkgLen = uint32(len(mesJsonData))
	binary.BigEndian.PutUint32(transfer.Buf[0:4], pkgLen)

	_, err = transfer.Conn.Write(transfer.Buf[0:4])
	if err != nil {
		fmt.Println("conn.Write header err=", err)
		return
	}

	// 发送消息体本身
	_, err = transfer.Conn.Write(mesJsonData)

	if err != nil {
		fmt.Println("conn.Write body err=", err)
		return
	}
	return nil
}
