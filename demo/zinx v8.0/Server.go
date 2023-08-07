package main

import (
	"zinx/znet"
	"zinx/znet/routers"
)

func main() {
	server := znet.NewServer()

	server.AddRouter(0, &routers.PingRouter{})
	server.AddRouter(1, &routers.FullDeliveryRouter{})

	server.Serve()
}
