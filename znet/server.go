package znet

import (
	"errors"
	"fmt"
	"net"
	"zinx/ziface"
)

type Server struct {
	Name      string
	IpVersion string
	Ip        string
	Port      int
	Conns     []Connection
	Router    ziface.IRouter
}

func NewServer(ip string, port int) *Server {
	server := &Server{
		Ip:        ip,
		IpVersion: "tcp4",
		Port:      port,
		Conns:     make([]Connection, 5),
		Router:    nil,
	}

	return server
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

		connection := NewConnection(tcpConn, uid, s.Router)
		uid++

		go connection.Establish()
		s.Conns = append(s.Conns, *connection)
	}
}

// ============== 定义当前客户端链接的handle api ===========
func CallBackToClient(conn *net.TCPConn, data []byte, cnt int) error {
	//回显业务
	fmt.Println("[Conn Handle] CallBackToClient ... ")
	if _, err := conn.Write(data[:cnt]); err != nil {
		fmt.Println("write back buf err ", err)
		return errors.New("CallBackToClient error")
	}
	return nil
}

// 释放资源并关闭服务器
func (s *Server) Stop() {
	for _, connection := range s.Conns {
		connection.Close()
	}
}

// server上线
func (s *Server) Serve() {
	s.start()
}

func (s *Server) AddRouter(router ziface.IRouter) {
	s.Router = router

	fmt.Println("Add router succ.")
}
