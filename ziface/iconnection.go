package ziface

import "net"

type IConnection interface {
	// connection建立
	Establish()

	// connection关闭
	Close()

	// 获取connection的uid
	GetConnId() uint32

	GetTCPConnection() *net.TCPConn
}

// 回调函数
// type HandleFunc func(conn *net.TCPConn, data []byte, cnt int) error
