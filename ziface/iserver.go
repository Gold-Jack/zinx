package ziface

type IServer interface {
	// server启动
	Start()

	// server停止
	Stop()

	// server上线服务
	Serve()

	// 路由功能：给当前服务器注册一个路由业务方法，供客户端连接使用
	AddRouter(router IRouter)
}
