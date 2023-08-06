package znet

import (
	"fmt"
	"net"
	"time"
	"zinx/utils"
	"zinx/ziface"
)

type Server struct {
	Name       string
	IpVersion  string
	Ip         string
	Port       int
	MsgHandler ziface.IMessageHandler
}

func NewServer() *Server {
	// 先加载全局配置文件
	utils.GlobalObject.Reload()

	s := &Server{
		Name:       utils.GlobalObject.Name,
		IpVersion:  "tcp4",
		Ip:         utils.GlobalObject.Host,
		Port:       utils.GlobalObject.TcpPort,
		MsgHandler: NewMessageHandler(),
	}

	return s
}

func (s *Server) start() {
	tcpAddr, err := net.ResolveTCPAddr(s.IpVersion, fmt.Sprintf("%s:%d", s.Ip, s.Port))
	if err != nil {
		fmt.Println("resolve tcp error.")
		return
	}
	listener, err := net.ListenTCP(s.IpVersion, tcpAddr)
	if err != nil {
		fmt.Println("listen tcp error.")
		return
	}

	fmt.Println("Server online.")
	var uid uint32 = 0

	for {
		tcpConn, err := listener.AcceptTCP()
		if err != nil {
			fmt.Println("accept tcp error.")
			continue
		}

		connection := NewConnection(tcpConn, uid, s.MsgHandler)
		uid++

		go connection.Establish()
		// s.Conns = append(s.Conns, *connection)
	}
}

// ============== 定义当前客户端链接的handle api ===========
// func CallBackToClient(conn *net.TCPConn, data []byte, cnt int) error {
// 	//回显业务
// 	fmt.Println("[Conn Handle] CallBackToClient ... ")
// 	if _, err := conn.Write(data[:cnt]); err != nil {
// 		fmt.Println("write back buf err ", err)
// 		return errors.New("CallBackToClient error")
// 	}
// 	return nil
// }

// 释放资源并关闭服务器
func (s *Server) Stop() {
	// for _, connection := range s.Conns {
	// 	connection.Close()
	// }
	fmt.Println("[Zinx STOP] Zinx server, name:", s.Name)
}

// server上线
func (s *Server) Serve() {
	// 这里打印的log主要是检测全局配置文件是否加载成功
	fmt.Printf("[Zinx INFO] Server name: %s,listenner at IP: %s, Port %d is starting\n", s.Name, s.Ip, s.Port)
	fmt.Printf("[Zinx INFO] Server version: %s, MaxConn: %d,  MaxPacketSize: %d\n",
		utils.GlobalObject.Version,
		utils.GlobalObject.MaxConn,
		utils.GlobalObject.MaxPacketSize)

	s.start()

	for {
		time.Sleep(10 * time.Second)
	}
}

func (s *Server) AddRouter(msgId uint32, router ziface.IRouter) {
	s.MsgHandler.AddRouter(msgId, router)
	// fmt.Println("Add router succ.")
}
