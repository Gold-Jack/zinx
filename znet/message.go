package znet

type Message struct {
	id      uint32 // 消息的ID
	dataLen uint32 // 消息的长度
	data    []byte // 消息的内容
}

func NewMessagePacket(id uint32, data []byte) *Message {
	return &Message{
		id:      id,
		data:    data,
		dataLen: uint32(len(data)),
	}
}

func (m *Message) GetDataLen() uint32 {
	return m.dataLen
}

func (m *Message) GetMsgId() uint32 {
	return m.id
}

func (m *Message) GetData() []byte {
	return m.data
}

func (m *Message) SetData(data []byte) {
	m.data = data
}

func (m *Message) SetMsgId(id uint32) {
	m.id = id
}

func (m *Message) SetDataLen(dataLen uint32) {
	m.dataLen = dataLen
}
