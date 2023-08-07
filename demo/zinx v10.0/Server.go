package main

import (
	"fmt"
	"zinx/ziface"
	"zinx/znet"
	"zinx/znet/routers"
)

func DoConnectionBegin(conn ziface.IConnection) {
	// fmt.Println("正在连接服务器")
	conn.SendMsg(0, []byte("正在连接服务器"))

	//=============设置两个链接属性，在连接创建之后===========
	conn.SetProperty("Name", "Gold_jack")
	conn.SetProperty("Home", "https://github.com/Gold-Jack/zinx")
	//===================================================
}

func DoConnectionEnd(conn ziface.IConnection) {
	// fmt.Println("与服务器断开连接")
	// conn.SendBuffMsg(0, []byte("与服务器断开连接"))
	//============在连接销毁之前，查询conn的Name，Home属性=====
	if name, err := conn.GetProperty("Name"); err == nil {
		fmt.Println("Conn Property Name = ", name)
	}

	if home, err := conn.GetProperty("Home"); err == nil {
		fmt.Println("Conn Property Home = ", home)
	}
	//===================================================
}

func main() {
	server := znet.NewServer()

	server.SetOnConnectionEstablish(DoConnectionBegin)
	server.SetOnConnectionStop(DoConnectionEnd)

	server.AddRouter(0, &routers.PingRouter{})
	server.AddRouter(1, &routers.FullDeliveryRouter{})

	server.Serve()
}
