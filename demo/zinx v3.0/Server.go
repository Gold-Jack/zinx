package main

import "zinx/znet"

func main() {
	ip, port := "127.0.0.1", 8888
	server := znet.NewServer(ip, port)
	server.AddRouter(&znet.PingRouter{})

	server.Serve()
}
