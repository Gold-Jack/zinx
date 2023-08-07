package main

import (
	"fmt"
	"io"
	"net"
	"time"
	"zinx/znet"
)

/*
模拟客户端
*/
func main() {

	fmt.Println("Client Test ... start")
	//3秒之后发起测试请求，给服务端开启服务的机会
	time.Sleep(3 * time.Second)

	conn, err := net.Dial("tcp4", "127.0.0.1:8888")
	if err != nil {
		fmt.Println("client start err, exit!")
		return
	}

	var uid uint32 = 0

	for {
		//发封包message消息
		dp := znet.NewDataPack()
		msg, _ := dp.Pack(znet.NewMessagePacket(uid, []byte("Zinx V0.5 Client Test Message")))
		uid++

		_, err := conn.Write(msg)
		if err != nil {
			fmt.Println("write error err ", err)
			return
		}

		//先读出流中的head部分
		headData := make([]byte, dp.GetHeadLen())
		_, err = io.ReadFull(conn, headData) //ReadFull 会把msg填充满为止
		if err != nil {
			fmt.Println("read head error")
			break
		}
		//将headData字节流 拆包到msg中
		msgHead, err := dp.Unpack(headData)
		if err != nil {
			fmt.Println("server unpack err:", err)
			return
		}

		if msgHead.GetDataLen() > 0 {
			//msg 是有data数据的，需要再次读取data数据
			msg := msgHead.(*znet.Message)
			msg.SetData(make([]byte, msg.GetDataLen()))

			//根据dataLen从io中读取字节流
			_, err := io.ReadFull(conn, msg.GetData())
			if err != nil {
				fmt.Println("server unpack data err:", err)
				return
			}

			// fmt.Println("==> Recv Msg: ID=", msg.GetMsgId(), ", len=", msg.GetDataLen(), ", data=", string(msg.GetData()))
			fmt.Printf("==> Recv Msg: Id=%d, len=%d, data=%v", msg.GetMsgId(), msg.GetDataLen(), string(msg.GetData()))
		}

		time.Sleep(1 * time.Second)
	}
}
