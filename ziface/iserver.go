package ziface

type IServer interface {
	// server启动
	Start()

	// server停止
	Stop()

	// server上线服务
	Serve()

	// 路由功能：给当前服务器注册一个路由业务方法，供客户端连接使用
	AddRouter(msgId uint32, router IRouter)

	// 得到连接管理
	GetConnectionManager() IConnectionManager

	// 设置Server的连接创建时的Hook函数
	SetOnConnectionEstablish(func(IConnection))

	// 设置Server的连接停止时的Hook函数
	SetOnConnectionStop(func(IConnection))

	// 调用连接创建时的Hook函数
	CallOnConnectionEstablish(conn IConnection)

	// 调用连接停止时的Hook函数
	CallOnConnectionStop(conn IConnection)
}
