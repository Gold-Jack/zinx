package znet

import (
	"errors"
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

	// const MAX_READ_BYTE = 512
	for {
		// 创建拆解包实例对象
		dp := NewDataPack()

		// 读取客户端的Message head
		header := make([]byte, dp.GetHeadLen())
		if _, err := io.ReadFull(c.GetTCPConnection(), header); err == io.EOF {
			fmt.Println("Client offline")
			c.ExitBuffChan <- true
			return
		} else if err != nil {
			fmt.Println("Read message head error", err)
			c.ExitBuffChan <- true
			continue
		}

		msg, err := dp.Unpack(header)
		if err != nil {
			fmt.Println("unpack error.")
			c.ExitBuffChan <- true
			continue
		}

		var data []byte
		if msg.GetDataLen() > 0 {
			data = make([]byte, msg.GetDataLen())
			if _, err := io.ReadFull(c.GetTCPConnection(), data); err != nil {
				fmt.Println("read msg data error")
				c.ExitBuffChan <- true
				continue
			}
		}
		msg.SetData(data)

		req := &Request{
			conn: c,
			msg:  msg,
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

func (c *Connection) SendMsg(msgId uint32, data []byte) error {
	if c.IsClosed {
		return errors.New("connection closed when send data")
	}

	// 将data封包，并发送给客户端
	dp := NewDataPack()
	msg, err := dp.Pack(NewMessagePacket(msgId, data))
	if err != nil {
		fmt.Println("Pack error msgId:", msgId)
		return errors.New("pack msg error")
	}

	// 写回客户端
	if _, err := c.Conn.Write(msg); err != nil {
		fmt.Println("Write msg error, msgId:", msgId)
		c.ExitBuffChan <- true
		return errors.New("write msg error")
	}

	return nil
}
