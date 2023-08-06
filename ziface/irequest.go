package ziface

type IRequest interface {
	// 获取请求连接的信息
	GetConnection() IConnection

	// 获取请求消息的数据
	GetData() []byte
}
