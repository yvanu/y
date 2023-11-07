package main

import (
	"fmt"
	"go-for-test/y"
	"net"
	"time"
)

type Ping struct {
}

func NewPing() y.Handler {
	return &Ping{}
}

func (p *Ping) HandleFunc(conn net.Conn, msg *y.Msg) *y.Msg {
	fmt.Println(msg)
	if y.IsNeedReply(msg.Typ) {
		conn.Write(y.Pack(y.NewMsg(y.UndoReply(msg.Typ), msg.Id, []byte("ping reply from server"))))
	}
	return nil
}

func (p *Ping) Biz(conn net.Conn, data any) {
	conn.Write(y.Pack(y.NewMsg(y.AddReply(y.PingPacket), y.GenId(), []byte("ping from server"))))
}

func main() {

	s := y.NewServer("30010")
	s.RegisterHandler(y.PingPacket, NewPing())
	go s.Serve()
	time.Sleep(10 * time.Second)
	s.MultiCast([]string{"shanghai"}, nil, y.PingPacket)

}
