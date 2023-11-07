package main

import (
	"fmt"
	"go-for-test/y"
	"net"
)

type Ping struct {
}

func NewPing() y.Handler {
	return &Ping{}
}

func (p *Ping) HandleFunc(conn net.Conn, msg *y.Msg) *y.Msg {
	fmt.Println(msg)
	if y.IsNeedReply(msg.Typ) {
		conn.Write(y.Pack(y.NewMsg(y.UndoReply(msg.Typ), msg.Id, []byte("ping reply from client"))))
	}
	return nil
}

func (p *Ping) Biz(conn net.Conn, data any) {
	conn.Write(y.Pack(y.NewMsg(y.AddReply(y.PingPacket), y.GenId(), []byte("ping from client"))))
}

func main() {
	client := y.NewClient("30010")
	client.RegisterHandler(y.PingPacket, NewPing())
	go client.Start()
	client.MultiCast(y.PingPacket, nil)
}
