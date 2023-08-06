package znet

import (
	"bytes"
	"encoding/binary"
	"errors"
	"zinx/utils"
	"zinx/ziface"
)

// 封包、拆包 实例类
type DataPack struct{}

func NewDataPack() *DataPack {
	return &DataPack{}
}

func (dp *DataPack) GetDefaultHeaderLen() uint32 {
	//Id uint32(4字节) +  DataLen uint32(4字节)
	return 8
}

// 封包（压缩数据）
func (dp *DataPack) Pack(msg ziface.IMessage) ([]byte, error) {
	dataBuff := bytes.NewBuffer([]byte{})

	// 写dataLen
	if err := binary.Write(dataBuff, binary.LittleEndian, msg.GetDataLen()); err != nil {
		return nil, err
	}

	// 写id
	if err := binary.Write(dataBuff, binary.LittleEndian, msg.GetMsgId()); err != nil {
		return nil, err
	}

	// 写data
	if err := binary.Write(dataBuff, binary.LittleEndian, msg.GetData()); err != nil {
		return nil, err
	}

	return dataBuff.Bytes(), nil
}

// 拆包（解压数据）
func (dp *DataPack) Unpack(binaryData []byte) (ziface.IMessage, error) {
	dataBuff := bytes.NewReader(binaryData)

	msg := &Message{}

	// 读dataLen
	if err := binary.Read(dataBuff, binary.LittleEndian, &msg.dataLen); err != nil {
		return nil, err
	}

	// 读id
	if err := binary.Read(dataBuff, binary.LittleEndian, &msg.id); err != nil {
		return nil, err
	}

	// 判断dataLen的长度是否超过了允许的最大包长度
	if utils.GlobalObject.MaxPacketSize > 0 && msg.dataLen > utils.GlobalObject.MaxPacketSize {
		return nil, errors.New("too large msg data received")
	}

	return msg, nil
}
