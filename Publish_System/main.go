package main

import (
	"log"
	"Publish_System/handler"
	"Publish_System/tool"
)

var config *tool.Config

func init()  {
	//加载配置文件
	config = tool.LoadConfig()
	if config == nil {
		return
	}

	//启动http服务来接收hook的请求
	go tool.NewHttpServer(config)
}

func main() {
	//创建微服务服务端

	s := tool.NewServer(config.ServerAddress,config.ConsulAddress,config.BasePath)
	err := s.RegisterName(config.ServerName,new(handler.ReceiveHookService), "")
	if err != nil {
		log.Fatalf("init server failed:%s",err)
	}

	s.Serve("tcp", config.ServerAddress)
}
