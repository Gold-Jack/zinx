package main

import (
	"zinx/znet"
)

func main() {
	server := znet.NewServer("127.0.0.1", 8888)
	server.Serve()
}
