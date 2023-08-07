package ziface

type IMessageHandler interface {
	// 以非阻塞方法处理消息
	DoMessageHandle(request IRequest)

	// 为消息添加具体的处理逻辑
	AddRouter(msgId uint32, router IRouter)

	// 开启worker工作池
	StartWorkerPool()

	// 将消息交给TaskQueue，由worker处理
	SendMsgToTaskQueue(request IRequest)
}
