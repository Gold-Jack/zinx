package routers

import (
	"fmt"
	"zinx/ziface"
)

type FullDeliveryRouter struct {
	BaseRouter
}

func (fdr *FullDeliveryRouter) PreHandle(request ziface.IRequest) {

}

func (fdr *FullDeliveryRouter) Handle(request ziface.IRequest) {
	fmt.Printf("Recv from client: {msgId=%d, data=%s}, sending...\n", request.GetMsgId(), request.GetData())
	err := request.GetConnection().SendMsg(request.GetMsgId(), request.GetData())
	if err != nil {
		fmt.Println(err)
	}
}

func (fdr *FullDeliveryRouter) PostHandle(request ziface.IRequest) {

}
