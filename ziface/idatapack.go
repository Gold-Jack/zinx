package ziface

type IDataPack interface {
	GetHeadLen() uint32                // 获取包header的长度
	Pack(msg IMessage) ([]byte, error) // 封包方法
	Unpack([]byte) (IMessage, error)
}
