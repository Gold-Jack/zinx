package znet

import (
	"fmt"
	"io"
	"net"
	"zinx/ziface"
)

type Connection struct {
	Conn         *net.TCPConn
	connId       uint32
	IsClosed     bool
	Router       ziface.IRouter // 该链接的处理方法
	ExitBuffChan chan bool      // 告知该connection已经停止的bufferChannel
}

func NewConnection(conn *net.TCPConn, connId uint32, router ziface.IRouter) *Connection {
	connection := &Connection{
		Conn:         conn,
		connId:       connId,
		IsClosed:     false,
		Router:       router,
		ExitBuffChan: make(chan bool, 1),
	}

	return connection
}

// 建立连接
func (c *Connection) Establish() {
	go c.StartDataReader()

	for {
		select {
		case <-c.ExitBuffChan:
			return
		}
	}
}

func (c *Connection) StartDataReader() {
	fmt.Println("DataReader starting...")
	defer fmt.Println("DataReader closing...")
	defer c.Close()

	const MAX_READ_BYTE = 512
	for {
		buf := make([]byte, MAX_READ_BYTE)
		_, err := c.Conn.Read(buf)
		if err == io.EOF {
			fmt.Println("client wants to close connection.")
			c.ExitBuffChan <- true
			return
		}
		if err != nil {
			fmt.Println("conn read error.")
			continue
		}

		// if err := c.CallBackFuncApi(c.Conn, buf, cnt); err != nil {
		// 	fmt.Println("callback error.")
		// 	c.ExitBuffChan <- true
		// 	return
		// }

		req := &Request{
			conn: c,
			data: buf,
		}

		go func(request ziface.IRequest) {
			// 执行注册路由的Handle方法集
			c.Router.PreHandle(request)
			c.Router.Handle(request)
			c.Router.PostHandle(request)
		}(req)
	}
}

// 释放资源并关闭连接
func (c *Connection) Close() {
	if c.IsClosed {
		return
	}

	c.IsClosed = true
	c.Conn.Close()
	c.ExitBuffChan <- true
	close(c.ExitBuffChan)
}

func (c *Connection) GetConnId() uint32 {
	return c.connId
}

func (c *Connection) GetTCPConnection() *net.TCPConn {
	return c.Conn
}
