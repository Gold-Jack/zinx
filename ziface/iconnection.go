package ziface

import "net"

type IConnection interface {
	// connection建立
	Establish()

	// connection关闭
	Close()

	// 获取connection的uid
	GetConnId() uint32

	// 获取Connection中的TCPConn连接实例
	GetTCPConnection() *net.TCPConn

	// 直接将Message发给远程的TCP客户端
	SendMsg(msgId uint32, data []byte) error

	// 带缓冲的发送方式
	SendBuffMsg(msgId uint32, data []byte) error
}

// 回调函数
// type HandleFunc func(conn *net.TCPConn, data []byte, cnt int) error
