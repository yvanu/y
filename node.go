package y

import "net"

type Node struct {
	probe string

	Conn net.Conn
}

func NewNode(probe string, conn net.Conn) *Node {
	return &Node{
		probe: probe,
		Conn:  conn,
	}
}
