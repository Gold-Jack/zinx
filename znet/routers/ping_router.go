package routers

import (
	"fmt"
	"zinx/ziface"
)

type PingRouter struct {
	BaseRouter
}

const PING_MSG string = "ping...ping...ping..."

func (pr *PingRouter) PreHandle(request ziface.IRequest) {
	// fmt.Println("Call Router PreHandle")

}

func (pr *PingRouter) Handle(request ziface.IRequest) {
	fmt.Println("Call PingRouter Handle")
	fmt.Printf("Recv from client: {msgId=%d, data=%s}, sending...\n", request.GetMsgId(), request.GetData())

	err := request.GetConnection().SendMsg(request.GetMsgId(), []byte(PING_MSG))
	if err != nil {
		fmt.Println(err)
	}
}

func (pr *PingRouter) PostHandle(request ziface.IRequest) {
	// fmt.Println("Call Router PostHandle")
}
