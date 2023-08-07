package main

import (
	"zinx/ziface"
	"zinx/znet"
	"zinx/znet/routers"
)

func DoConnectionBegin(conn ziface.IConnection) {
	// fmt.Println("正在连接服务器")
	conn.SendBuffMsg(0, []byte("正在连接服务器"))
}

func DoConnectionEnd(conn ziface.IConnection) {
	// fmt.Println("与服务器断开连接")
	// conn.SendBuffMsg(0, []byte("与服务器断开连接"))
}

func main() {
	server := znet.NewServer()

	server.SetOnConnectionEstablish(DoConnectionBegin)
	server.SetOnConnectionStop(DoConnectionEnd)

	server.AddRouter(0, &routers.PingRouter{})
	server.AddRouter(1, &routers.FullDeliveryRouter{})

	server.Serve()
}
