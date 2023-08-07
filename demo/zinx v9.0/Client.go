package main

import (
	"fmt"
	"io"
	"math/rand"
	"net"
	"time"
	"zinx/znet"
)

func main() {
	for i := 0; i < 1; i++ {
		go ClientConnect(uint32(i + 1))
	}

	time.Sleep(100 * time.Second)
}

func ClientConnect(clientId uint32) {
	network, address := "tcp4", "127.0.0.1:8888"
	TcpConn, err := net.Dial(network, address)
	if err != nil {
		fmt.Println("Tcp connection failed.")
		return
	}

	dp := znet.NewDataPack()
	for {
		// fmt.Printf("Client-%d\n")
		msgId := rand.Intn(2)
		msg, _ := dp.Pack(znet.NewMessagePacket(uint32(msgId), []byte("Zinx V9.0 test message.")))

		_, err := TcpConn.Write(msg)
		if err != nil {
			fmt.Println("tcpconn write message error.")
			return
		}

		// 先读出head部分
		headData := make([]byte, dp.GetDefaultHeaderLen())
		_, err = io.ReadFull(TcpConn, headData)
		if err != nil {
			fmt.Println("header read error.")
			return
		}

		msgHead, err := dp.Unpack(headData)
		if err != nil {
			fmt.Println("unpack header error.")
		}

		if msgHead.GetDataLen() > 0 {
			msg := msgHead.(*znet.Message)
			msg.SetData(make([]byte, msgHead.GetDataLen()))

			_, err = io.ReadFull(TcpConn, msg.GetData())
			if err != nil {
				fmt.Println("data read error.")
				return
			}

			fmt.Printf("[Client %v] ==> Recv Msg: Id=%d, len=%d, data=%v\n", clientId, msg.GetMsgId(), msg.GetDataLen(), string(msg.GetData()))
		}

		time.Sleep(1 * time.Second)
	}
}
