package y

import (
	"fmt"
	"io"
	"net"
)

// typ
// 首位为1表示要回复，为0表示不需要
// 0x80000001 表示心跳包
// 0x80000002 表示新入node 注册包

const (
	BeatPacket = iota + 0x80000001
	RegisPacket

	PingPacket
)

type Msg struct {
	Id      uint32
	Typ     uint32 // 用户自定义
	DataLen uint32
	Data    []byte
}

func (m *Msg) String() string {
	return fmt.Sprintf("id: %d\ntype: %d\ndatalen: %d\ndata: %s\n", m.Id, m.Typ, m.DataLen, string(m.Data))
}

func NewMsg(typ, id uint32, data []byte) *Msg {
	m := &Msg{}
	if id == 0 {
		id = genId()
	}
	m.Id = id
	m.Typ = typ
	m.Data = data
	m.DataLen = uint32(len(data))
	return m
}

func ReadMsg(conn net.Conn) (*Msg, error) {
	msg := &Msg{}
	header := make([]byte, 12)
	_, err := conn.Read(header)
	if err != nil {
		//if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
		//	fmt.Println("读取超时，关闭链接")
		//	return nil, err
		//}
		//if err == io.EOF {
		//	fmt.Println("对端关闭链接")
		//	return nil, err
		//}
		return nil, err
	}
	msg, err = DefaultDataPack.UnPack(header)
	if err != nil {
		fmt.Println("拆包错误")
	}

	if msg.DataLen > 0 {
		msg.Data = make([]byte, msg.DataLen)
		_, err = io.ReadFull(conn, msg.Data)
		if err != nil {
			fmt.Println("读取data失败")
			panic(err)
		}
	}
	return msg, nil
}
