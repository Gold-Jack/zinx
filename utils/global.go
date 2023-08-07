package utils

import (
	"encoding/json"
	"os"
	"zinx/ziface"
)

type GlobalObj struct {
	/*
		Server
	*/
	TcpServer ziface.IServer // 当前Zinx的全局Server对象
	Host      string         // 当前服务器主机IP
	TcpPort   int            // 当前服务器主机监听端口号
	Name      string         // 当前服务器名称

	/*
		Zinx
	*/
	Version            string // 当前Zinx版本号
	MaxPacketSize      uint32 // 当前数据包的最大值
	MaxConn            int    // 当前服务器主机允许的最大链接个数
	WorkerPoolSize     uint32 // 池中worker的数量
	MaxWorkerTaskLen   uint32 // Worker对应负责的任务队列最大任务存储数量
	MaxConnectionCount int    // 最大连接数
	MaxMsgBuffChanLen  int    // 最大带缓冲管道的byte数

	/*
		Config file path
	*/
	ConfigFilePath string
}

// 定义一个全局的变量
var GlobalObject *GlobalObj

func (g *GlobalObj) Reload() {
	data, err := os.ReadFile("config/application.json")
	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(data, &GlobalObject)
	if err != nil {
		panic(err)
	}
}

func init() {
	// 初始化GlobalObject变量，设置默认值
	GlobalObject = &GlobalObj{
		Name:               "ZinxServerApp",
		Version:            "V4.0",
		TcpPort:            8888,
		Host:               "0.0.0.0",
		MaxConn:            12000,
		MaxPacketSize:      4096,
		WorkerPoolSize:     10,
		MaxWorkerTaskLen:   1024,
		MaxConnectionCount: 3,
		MaxMsgBuffChanLen:  1024,
		ConfigFilePath:     "config/application.json",
	}

	GlobalObject.Reload()
}
