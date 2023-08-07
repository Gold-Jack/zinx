package znet

import (
	"errors"
	"fmt"
	"io"
	"net"
	"sync"
	"zinx/utils"
	"zinx/ziface"
)

type Connection struct {
	Server       ziface.IServer // connection连接的Server
	Conn         *net.TCPConn
	connId       uint32
	IsClosed     bool
	MsgHandler   ziface.IMessageHandler // 消息管理模块
	ExitBuffChan chan bool              // 告知该connection已经停止的bufferChannel
	msgChan      chan []byte            // 无缓冲管道，用于读、写两个goroutine之间的通信
	msgBuffChan  chan []byte            // 有缓冲管道

	/*
		用户信息
	*/
	property     map[string]interface{}
	propertyLock sync.RWMutex
}

func NewConnection(server ziface.IServer, conn *net.TCPConn, connId uint32, msgHandler ziface.IMessageHandler) *Connection {
	connection := &Connection{
		Server:       server,
		Conn:         conn,
		connId:       connId,
		IsClosed:     false,
		MsgHandler:   msgHandler,
		ExitBuffChan: make(chan bool, 1),
		msgChan:      make(chan []byte),
		msgBuffChan:  make(chan []byte, utils.GlobalObject.MaxMsgBuffChanLen),
		property:     make(map[string]interface{}),
		propertyLock: sync.RWMutex{},
	}

	return connection
}

// 建立连接
func (c *Connection) Establish() {
	// 开启用于从客户端读取数据的reader
	go c.StartReader()

	// 开启用于将数据写回客户端的writer
	go c.StartWriter()

	c.Server.CallOnConnectionEstablish(c)
	c.Server.GetConnectionManager().Add(c)

	for {
		select {
		case <-c.ExitBuffChan:
			// connection已经关闭
			fmt.Println("Client offline")
			return
		}
	}
}

func (c *Connection) StartReader() {
	fmt.Println("[Zinx ESTABLISH] Connection reader starting...")
	defer fmt.Println("[Zinx CLOSE] Connection reader closing...")
	defer c.Close()

	for {
		// 创建拆解包实例对象
		dp := NewDataPack()

		// 读取客户端的Message head
		headData := make([]byte, dp.GetDefaultHeaderLen())
		if _, err := io.ReadFull(c.GetTCPConnection(), headData); err == io.EOF {
			// 用户主动结束连接
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

		if utils.GlobalObject.WorkerPoolSize > 0 {
			// 如果已经启用了workerPool
			c.MsgHandler.SendMsgToTaskQueue(request)
		} else {
			// 使用消息管理模块执行DoMessageHandler操作
			go c.MsgHandler.DoMessageHandle(request)
		}
	}
}

func (c *Connection) StartWriter() {
	fmt.Println("[Zinx ESTABLISH] Connection writer starting...")
	defer fmt.Println("[Zinx CLOSE] Connection writer closing...")
	defer c.Close()

	for {
		select {
		case data := <-c.msgChan:
			// 无缓冲管道中，有数据要写给客户端
			if _, err := c.Conn.Write(data); err != nil {
				fmt.Println("send data error")
				return
			}
		case buffData := <-c.msgBuffChan:
			// 缓冲管道中，有数据需要写给客户端
			if _, err := c.Conn.Write(buffData); err != nil {
				fmt.Println("Send buff data error")
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

	// 执行连接关闭前的hook函数
	c.Server.CallOnConnectionStop(c)

	// 关闭socket连接
	c.Conn.Close()

	// 发送connection关闭通知
	c.ExitBuffChan <- true

	// 从manager中移除当前connection
	c.Server.GetConnectionManager().Remove(c)

	// 关闭所有管道
	close(c.ExitBuffChan)
	close(c.msgChan)
	fmt.Println("[Zinx CLOSE] Connection closed.")
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
	c.msgChan <- msg

	return nil
}

func (c *Connection) SendBuffMsg(msgId uint32, data []byte) error {
	if c.IsClosed {
		return errors.New("connection closed when send data")
	}

	dp := NewDataPack()
	msg, err := dp.Pack(NewMessagePacket(msgId, data))
	if err != nil {
		fmt.Println("Pack error msgId:", msgId)
		return errors.New("pack msg error")
	}

	c.msgBuffChan <- msg

	return nil
}

func (c *Connection) SetProperty(key string, value interface{}) {
	c.propertyLock.Lock()
	defer c.propertyLock.Unlock()

	c.property[key] = value
}

func (c *Connection) GetProperty(key string) (interface{}, error) {
	c.propertyLock.Lock()
	defer c.propertyLock.Unlock()

	if value, exist := c.property[key]; exist {
		return value, nil
	} else {
		return nil, errors.New("no property found")
	}
}

func (c *Connection) RemoveProperty(key string) {
	c.propertyLock.Lock()
	defer c.propertyLock.Unlock()

	delete(c.property, key)
}
