package y

import (
	"fmt"
	"net"
	"sync"
	"time"
)

var port = "30010"

type Server struct {
	Port       string
	HandlerMap map[uint32]Handler

	nodeManager NodeManager

	connMap map[string]*Node
	lock    sync.RWMutex
}

func NewServer(port string) *Server {
	return &Server{Port: port, HandlerMap: make(map[uint32]Handler), connMap: make(map[string]*Node)}
}

func (s *Server) RegisterHandler(typ uint32, handler Handler) {
	s.HandlerMap[typ] = handler
}

func (s *Server) Serve(ok chan struct{}) {

	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		panic(err)
	}
	defer lis.Close()

	ok <- struct{}{}
	for {
		conn, err := lis.Accept()
		if err != nil {
			return
		}
		go s.process(conn)
	}
}

func setTimeout(conn net.Conn, t int) {
	conn.SetDeadline(time.Now().Add(time.Duration(t) * time.Second))
}

func (s *Server) process(conn net.Conn) {
	defer conn.Close()
	for {
		setTimeout(conn, 10)
		msg, err := ReadMsg(conn)
		if err != nil {
			// 删除该conn
			break
		}
		// 对msg进行业务处理
		go s.handleMsg(conn, msg)
	}
}

func (s *Server) handleMsg(conn net.Conn, msg *Msg) {
	// 心跳包 直接返回ok
	if isBeat(msg.Typ) {
		fmt.Printf("收到心跳包 %s", msg.String())
		retMsg := Pack(NewMsg(UndoReply(BeatPacket), msg.Id, []byte("beat ok")))
		conn.Write(retMsg)
		return
	}

	// 注册包
	if isRegister(msg.Typ) {
		fmt.Printf("收到注册包 %s", msg.String())
		s.lock.Lock()
		//if s.nodeManager.HasNode(string(msg.Data)) {
		if _, exist := s.connMap[string(msg.Data)]; exist {
			conn.Write(Pack(NewMsg(UndoReply(RegisPacket), msg.Id, []byte("已存在该节点"))))
			s.lock.Unlock()
			return
		}
		//s.nodeManager.AddNode(string(msg.Data), NewNode(string(msg.Data), conn))
		s.connMap[string(msg.Data)] = NewNode(string(msg.Data), conn)
		conn.Write(Pack(NewMsg(UndoReply(RegisPacket), msg.Id, []byte("register ok"))))
		s.lock.Unlock()

		return
	}
	fmt.Printf("收到业务包 %s", msg.String())
	// 其他业务包
	handler, exist := s.HandlerMap[msg.Typ]
	if !exist {
		fmt.Println(msg)
		panic("未注册的业务！")
	}
	handler.HandleFunc(conn, msg)
}

// MultiCast 把请求发给多个probes
func (s *Server) MultiCast(probes []string, data any, typ uint32) {
	for _, probe := range probes {
		s.lock.RLock()
		node, exist := s.connMap[probe]
		s.lock.RUnlock()
		if !exist {
			fmt.Printf("probe: %s 未注册", probe)
			continue
		}
		//node, err := s.nodeManager.GetNode(probe)
		//if err != nil{
		//	fmt.Println(err)
		//	continue
		//}
		s.HandlerMap[typ].Biz(node.Conn, data)
	}
}
