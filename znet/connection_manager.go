package znet

import (
	"errors"
	"fmt"
	"sync"
	"zinx/utils"
	"zinx/ziface"
)

type ConnectionManager struct {
	connections map[uint32]ziface.IConnection // 管理的连接信息
	RWLock      sync.RWMutex                  // 读写锁
}

func NewConnectionManager() *ConnectionManager {
	ConnMgr := &ConnectionManager{
		connections: make(map[uint32]ziface.IConnection, utils.GlobalObject.MaxConnectionCount),
		RWLock:      sync.RWMutex{},
	}

	return ConnMgr
}

func (cm *ConnectionManager) Add(conn ziface.IConnection) bool {
	cm.RWLock.Lock()
	if cm.Count() >= utils.GlobalObject.MaxConnectionCount {
		defer conn.Close()
		fmt.Println("Connection is full, please try again.")
		return false
	}

	if _, exist := cm.connections[conn.GetConnId()]; exist {
		defer conn.Close()
		fmt.Println("ConnId already exists, please try again.")
		return false
	}

	// 添加conn到manager.connections
	cm.connections[conn.GetConnId()] = conn
	cm.RWLock.Unlock()
	fmt.Printf("Add conn{id=%d} to ConnectionManager succ.\n", conn.GetConnId())
	return true
}

func (cm *ConnectionManager) Remove(conn ziface.IConnection) bool {
	if _, exist := cm.connections[conn.GetConnId()]; !exist {
		fmt.Printf("Conn{id=%d} does not exist, remove failed.\n", conn.GetConnId())
		return false
	}

	cm.RWLock.Lock()
	delete(cm.connections, conn.GetConnId())
	cm.RWLock.Unlock()
	fmt.Printf("Remove conn{id=%d} from ConnectionManager succ.\n", conn.GetConnId())
	return true
}

func (cm *ConnectionManager) Get(connId uint32) (ziface.IConnection, error) {
	// var conn ziface.IConnection
	cm.RWLock.Lock()
	conn, exist := cm.connections[connId]
	cm.RWLock.Unlock()
	if exist {
		return conn, nil
	} else {
		return nil, errors.New("conn not found, get failed")
	}
}

func (cm *ConnectionManager) Count() int {
	return len(cm.connections)
}

func (cm *ConnectionManager) Clear() {
	cm.RWLock.Lock()
	cnt := cm.Count()
	for connId, conn := range cm.connections {
		// 先停止conn的业务
		conn.Close()
		// 再删除
		delete(cm.connections, connId)
	}
	cm.RWLock.Unlock()
	fmt.Printf("Clear all connections succ, len=%d.\n", cnt)
}
