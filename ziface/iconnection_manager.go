package ziface

type IConnectionManager interface {
	Add(conn IConnection) bool              // 添加一个连接
	Remove(conn IConnection) bool           // 移除一个连接
	Get(connId uint32) (IConnection, error) // 通过connId获取连接
	Count() int                             // 获取当前连接数量
	Clear()                                 // 删除并停止所有连接
}
