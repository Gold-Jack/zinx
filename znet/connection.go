package znet

import (
	"errors"
	"fmt"
	"io"
	"net"
	"zinx/ziface"
)

type Connection struct {
	Conn          *net.TCPConn
	connId        uint32
	IsClosed      bool
	MsgHandler    ziface.IMessageHandler // 消息管理模块
	ExitBuffChan  chan bool              // 告知该connection已经停止的bufferChannel
	readWriteChan chan []byte            // 无缓冲管道，用于读、写两个goroutine之间的通信
}

func NewConnection(conn *net.TCPConn, connId uint32, msgHandler ziface.IMessageHandler) *Connection {
	connection := &Connection{
		Conn:          conn,
		connId:        connId,
		IsClosed:      false,
		MsgHandler:    msgHandler,
		ExitBuffChan:  make(chan bool, 1),
		readWriteChan: make(chan []byte),
	}

	return connection
}

// 建立连接
func (c *Connection) Establish() {
	// 开启用于从客户端读取数据的reader
	go c.StartReader()

	// 开启用于将数据写回客户端的writer
	go c.StartWriter()

	for {
		select {
		case <-c.ExitBuffChan:
			// connection已经关闭
			return
		}
	}
}

func (c *Connection) StartReader() {
	fmt.Println("[Zinx ESTABLISH] Reader starting...")
	defer fmt.Println("[Zinx ESTABLISH] Reader closing...")
	defer c.Close()

	for {
		// 创建拆解包实例对象
		dp := NewDataPack()

		// 读取客户端的Message head
		headData := make([]byte, dp.GetDefaultHeaderLen())
		if _, err := io.ReadFull(c.GetTCPConnection(), headData); err == io.EOF {
			fmt.Println("Client offline")
			c.ExitBuffChan <- true
			return
		} else if err != nil {
			fmt.Println("Read message head error", err)
			c.ExitBuffChan <- true
			continue
		}

		msgHead, err := dp.Unpack(headData)
		if err != nil {
			fmt.Println("unpack error.")
			c.ExitBuffChan <- true
			continue
		}

		var msg *Message
		if msgHead.GetDataLen() > 0 {
			msg = msgHead.(*Message)
			msg.SetData(make([]byte, msgHead.GetDataLen()))
			if _, err := io.ReadFull(c.GetTCPConnection(), msg.GetData()); err != nil {
				fmt.Println("read msg data error")
				c.ExitBuffChan <- true
				continue
			}
		}

		request := &Request{
			conn: c,
			msg:  msg,
		}

		// 使用消息管理模块执行DoMessageHandler操作
		go c.MsgHandler.DoMessageHandle(request)
	}
}

func (c *Connection) StartWriter() {
	fmt.Println("[Zinx ESTABLISH] Writer starting...")
	defer fmt.Println("[Zinx ESTABLISH] Writer closing...")
	defer c.Close()

	for {
		select {
		case data := <-c.readWriteChan:
			// 有数据要写给客户端
			if _, err := c.Conn.Write(data); err != nil {
				fmt.Println("send data error")
				return
			}
		case <-c.ExitBuffChan:
			// connection已经关闭
			return
		}
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
	c.readWriteChan <- msg

	return nil
}
