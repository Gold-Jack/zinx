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
	ConnMgr    ziface.IConnectionManager

	/*
		两个Hook函数原型
	*/
	OnConnEstablish func(conn ziface.IConnection)
	OnConnStop      func(conn ziface.IConnection)
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
		ConnMgr:    NewConnectionManager(),
	}

	return s
}

func (s *Server) Start() {
	s.MsgHandler.StartWorkerPool()

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

	fmt.Println("[Zinx START] Server online.")

	var uid uint32 = 0
	for {
		tcpConn, err := listener.AcceptTCP()
		if err != nil {
			fmt.Println("accept tcp error.")
			continue
		}

		connection := NewConnection(s, tcpConn, uid, s.MsgHandler)
		uid++

		go connection.Establish()
	}
}

// 释放资源并关闭服务器
func (s *Server) Stop() {
	defer fmt.Printf("[Zinx STOP] Zinx server{name=%s} stopped.\n", s.Name)
	// 清除manager中的所有连接
	s.ConnMgr.Clear()
}

// server上线
func (s *Server) Serve() {
	// 这里打印Zinx的基本信息
	// 检测全局配置文件是否加载成功
	fmt.Printf("[Zinx INFO] Server name: %s,listenner at IP: %s, Port %d is starting\n",
		s.Name,
		s.Ip,
		s.Port)
	fmt.Printf("[Zinx INFO] Server version: %s, MaxConn: %d,  MaxPacketSize: %d\n",
		utils.GlobalObject.Version,
		utils.GlobalObject.MaxConn,
		utils.GlobalObject.MaxPacketSize)

	s.Start()

	for {
		time.Sleep(10 * time.Second)
	}
}

func (s *Server) AddRouter(msgId uint32, router ziface.IRouter) {
	s.MsgHandler.AddRouter(msgId, router)
}

func (s *Server) GetConnectionManager() ziface.IConnectionManager {
	return s.ConnMgr
}

func (s *Server) SetOnConnectionEstablish(hookFunc func(conn ziface.IConnection)) {
	s.OnConnEstablish = hookFunc
}

func (s *Server) SetOnConnectionStop(hookFun func(conn ziface.IConnection)) {
	s.OnConnStop = hookFun
}

func (s *Server) CallOnConnectionEstablish(conn ziface.IConnection) {
	if s.OnConnEstablish != nil {
		s.OnConnEstablish(conn)
	}
}

func (s *Server) CallOnConnectionStop(conn ziface.IConnection) {
	if s.OnConnStop != nil {
		s.OnConnStop(conn)
	}
}
