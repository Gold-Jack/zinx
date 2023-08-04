package ziface

type IServer interface {
	// server启动
	Start()

	// server停止
	Stop()

	// server上线服务
	Serve()
}
