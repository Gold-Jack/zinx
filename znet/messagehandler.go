package znet

import (
	"fmt"
	"zinx/ziface"
	"zinx/znet/routers"
)

type MessageHandler struct {
	Apis map[uint32]ziface.IRouter //存放每个MsgId 所对应的处理方法
	// WorkerPoolSize uint32
	// TaskQueue      []chan ziface.IRequest
}

func NewMessageHandler() *MessageHandler {
	return &MessageHandler{
		Apis: make(map[uint32]ziface.IRouter),
	}
}

func (mh *MessageHandler) DoMessageHandle(request ziface.IRequest) {
	router, exist := mh.Apis[request.GetMsgId()]
	if !exist {
		// 如果没有对应的Router，则默认为PingRouter
		router = &routers.PingRouter{}
	}

	// 调用Router的handle方法集
	router.PreHandle(request)
	router.Handle(request)
	router.PostHandle(request)
}

func (mh *MessageHandler) AddRouter(msgId uint32, router ziface.IRouter) {
	// 1. 判断当前msgId是否已经绑定了Router
	if _, exist := mh.Apis[msgId]; exist {
		panic(fmt.Sprintf("[Zinx WARN] Repeated router, msgId=%d\n", msgId))
	}

	// 2. 绑定msgId和router
	mh.Apis[msgId] = router
	fmt.Printf("[Zinx MOUNT] Add router succ, msgId=%d\n", msgId)
}
