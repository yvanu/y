package y

import "net"

type Handler interface {
	HandleFunc(conn net.Conn, msg *Msg) *Msg
	Biz(conn net.Conn, data any)
}
