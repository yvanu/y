package y

import (
	"errors"
	"fmt"
	"net"
	"sync"
)

type NodeManager struct {
	nodeMap map[string]*Node
	lock    sync.RWMutex
}

func NewNodeManager() *NodeManager {
	return &NodeManager{nodeMap: make(map[string]*Node)}
}

func (n *NodeManager) HasNode(name string) bool {
	_, exist := n.nodeMap[name]
	if exist {
		return false
	}
	return true
}

func (n *NodeManager) AddNode(name string, node *Node) {
	n.lock.Lock()
	defer n.lock.Unlock()
	n.nodeMap[name] = node

}

func (n *NodeManager) RemoveByNode(name string) {
	n.lock.Lock()
	defer n.lock.Unlock()
	delete(n.nodeMap, name)

}

func (n *NodeManager) RemoveByConn(conn net.Conn) {
	n.lock.Lock()
	defer n.lock.Unlock()
	for name, node := range n.nodeMap {
		if node.Conn == conn {
			delete(n.nodeMap, name)
		}
	}
}

func (n *NodeManager) GetNode(name string) (*Node, error) {
	n.lock.Lock()
	defer n.lock.Unlock()
	node, exist := n.nodeMap[name]
	if !exist {
		return nil, errors.New(fmt.Sprintf("该node: %s不存在", name))
	}
	return node, nil
}
