package znet

import (
	"errors"
	"fmt"
	"sync"
	"zinx/ziface"
)

type ConnManager struct {
	connections map[int32]ziface.IConnection
	lock        sync.RWMutex
}

func NewConnManager() *ConnManager {
	return &ConnManager{
		connections: make(map[int32]ziface.IConnection),
	}
}

func (cm *ConnManager) Add(connId int32, conn ziface.IConnection) {
	defer cm.lock.Unlock()

	cm.lock.Lock()
	cm.connections[connId] = conn
	fmt.Println("connection add to ConnManager successfully: conn num = ", cm.Len())
}

func (cm *ConnManager) Remove(connId int32) {
	defer cm.lock.Unlock()

	cm.lock.Lock()
	delete(cm.connections, connId)
	fmt.Println("connection Remove ConnID=", connId, " successfully: conn num = ", cm.Len())
}

func (cm *ConnManager) Get(connId int32) (ziface.IConnection, error) {
	defer cm.lock.RUnlock()

	cm.lock.RLock()
	if conn, ok := cm.connections[connId]; ok {
		return conn, nil
	} else {
		return nil, errors.New("connection not found")
	}
}

func (cm *ConnManager) Len() int {
	return len(cm.connections)
}

// 清空并关闭所有连接
func (cm *ConnManager) ClearAll() {}
