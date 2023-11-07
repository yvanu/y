package y

import "net"

type Handler interface {

	// 收包逻辑
	HandleFunc(conn net.Conn, msg *Msg) *Msg
	// 发包逻辑
	Biz(conn net.Conn, data any)
}
