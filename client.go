package y

import (
	"errors"
	"fmt"
	"net"
	"time"
)

type Client struct {
	serverPort string
	HandlerMap map[uint32]Handler
	conn       net.Conn
}

func NewClient(serverPort string) *Client {
	return &Client{serverPort: serverPort, HandlerMap: make(map[uint32]Handler)}

}

func (c *Client) RegisterHandler(typ uint32, handler Handler) {
	c.HandlerMap[typ] = handler
}

func sendBeat(conn net.Conn) {
	msg, _ := DefaultDataPack.Pack(NewMsg(BeatPacket, 0, nil))
	conn.Write(msg)
}

func keepalive(conn net.Conn) {
	for {
		sendBeat(conn)
		time.Sleep(3 * time.Second)
	}
}

func register(conn net.Conn) error {
	regisMsg := Pack(NewMsg(RegisPacket, 0, []byte("shanghai")))
	conn.Write(regisMsg)
	msg, err := ReadMsg(conn)
	if err != nil {
		return err
	}
	if string(msg.Data) != "register ok" {
		return errors.New("注册失败, " + string(msg.Data))
	}
	return nil

}

func (c *Client) Start() {
	conn, err := net.Dial("tcp", ":"+c.serverPort)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	// 先注册
	err = register(conn)

	if err != nil {
		panic(err)
	}

	go keepalive(conn)
	for {
		msg, _ := ReadMsg(conn)
		go c.handleMsg(conn, msg)
	}
}

func (c *Client) handleMsg(conn net.Conn, msg *Msg) {
	if isBeat(msg.Typ) {
		fmt.Println("收到心跳回包")
		return
	}

	// 其他业务包
	handler, exist := c.HandlerMap[msg.Typ]
	if !exist {
		fmt.Println(msg)
		panic("未注册的业务！")
	}
	handler.HandleFunc(conn, msg)
}

func (c *Client) MultiCast(typ uint32, data any) {
	handler, exist := c.HandlerMap[typ]
	if !exist {
		fmt.Printf("handler %d 未注册", typ)
		return
	}
	handler.Biz(c.conn, data)
}
