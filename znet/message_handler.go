package znet

import (
	"fmt"
	"zinx/utils"
	"zinx/ziface"
	"zinx/znet/routers"
)

type MessageHandler struct {
	Apis           map[uint32]ziface.IRouter //存放每个MsgId 所对应的处理方法
	WorkerPoolSize uint32                    //业务工作Worker的数量
	TaskQueue      []chan ziface.IRequest    //Worker负责取任务的消息队列
}

func NewMessageHandler() *MessageHandler {
	return &MessageHandler{
		Apis:           make(map[uint32]ziface.IRouter),
		WorkerPoolSize: utils.GlobalObject.WorkerPoolSize,
		TaskQueue:      make([]chan ziface.IRequest, utils.GlobalObject.WorkerPoolSize), // 一个worker绑定一个taskQueue
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

func (mh *MessageHandler) StartWorker(workerId int, taskQueue chan ziface.IRequest) {
	fmt.Printf("[Zinx START] MessageHandler Worker-%d starts.\n", workerId)
	// 绑定消息队列并监听
	for {
		select {
		case request := <-taskQueue:
			mh.DoMessageHandle(request)
		}
	}
}

func (mh *MessageHandler) StartWorkerPool() {
	// 启动workerPoolSize规定数量的worker
	for i := 0; i < int(mh.WorkerPoolSize); i++ {
		// 给当前worker开辟taskQueue的空间
		mh.TaskQueue[i] = make(chan ziface.IRequest, utils.GlobalObject.MaxWorkerTaskLen)
		go mh.StartWorker(i, mh.TaskQueue[i])
	}
}

func (mh *MessageHandler) SendMsgToTaskQueue(request ziface.IRequest) {
	// 轮询分配法则
	// 根据ConnId来分配当前的连接应该由哪个worker处理
	workerId := request.GetConnection().GetConnId() % mh.WorkerPoolSize
	fmt.Printf("Add Request{ConnId=%d, msgId=%d} to Worker{id=%d}",
		request.GetConnection().GetConnId(), request.GetMsgId(), workerId)
	mh.TaskQueue[workerId] <- request
}
