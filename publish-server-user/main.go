package main

import (
	"flag"
	"github.com/rcrowley/go-metrics"
	"github.com/smallnest/rpcx/server"
	"github.com/smallnest/rpcx/serverplugin"
	"log"
	"publish_server_user/datasource"
	"publish_server_user/service"
	"publish_server_user/tool"
	"time"
)

var (
	addr       = flag.String("addr", ":"+tool.ConfJson.PublishPort, "server address")
	consulAddr = flag.String("consulAddr", tool.ConfJson.ConsulAddr, "consul address")
	basePath   = flag.String("base", tool.ConfJson.BasePath, "prefix path")
)

func main() {
	defer datasource.GetDB().Close()
	datasource.CreateTable() //创建数据库表
	flag.Parse()
	s := server.NewServer()
	addRegistryPlugin(s)
	s.RegisterName("zzy.server.user", new(service.UserServices), "")
	//go service.GetNSQWarehouseProjectName()
	err := s.Serve("tcp", *addr)
	if err != nil {
		log.Fatal(err)
	}
}
func addRegistryPlugin(s *server.Server) {
	r := &serverplugin.ConsulRegisterPlugin{
		ServiceAddress: "tcp@" + *addr,
		ConsulServers:  []string{*consulAddr},
		BasePath:       *basePath,
		Metrics:        metrics.NewRegistry(),
		UpdateInterval: time.Minute,
	}
	err := r.Start()
	if err != nil {
		log.Fatal(err)
	}
	s.Plugins.Add(r)
}
