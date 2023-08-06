package znet

import "zinx/ziface"

type Request struct {
	conn ziface.IConnection // 客户端的连接
	msg  ziface.IMessage    // 客户端的请求消息
}

func (r *Request) GetConnection() ziface.IConnection {
	return r.conn
}

func (r *Request) GetData() []byte {
	return r.msg.GetData()
}

func (r *Request) GetMsgId() uint32 {
	return r.msg.GetMsgId()
}
