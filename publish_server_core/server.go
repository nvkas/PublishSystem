package main

import (
	"flag"
	"fmt"
	"github.com/rcrowley/go-metrics"
	"github.com/smallnest/rpcx/server"
	"github.com/smallnest/rpcx/serverplugin"
	"publish_server_core/services"
	"publish_server_core/tool"
	"time"
)

func main() {
	flag.Parse()
	s := server.NewServer()
	addRegistryPlugin(tool.RpcConfig.Addr,tool.RpcConfig.BasePath,s)
	s.RegisterName("server.core.allCommand", new(service.SettingServer), "")
	s.Serve("tcp", tool.RpcConfig.Addr)

}

func addRegistryPlugin(addr string,basePath string,s *server.Server) {
	r := &serverplugin.ConsulRegisterPlugin{
		ServiceAddress: "tcp@" + addr,
		ConsulServers:  []string{tool.RpcConfig.ConsulAddr},
		BasePath:       basePath,
		Metrics:        metrics.NewRegistry(),
		UpdateInterval: time.Minute,
	}

	err := r.Start()
	if err != nil {
		fmt.Println("Error:",err)
	}
	s.Plugins.Add(r)
}
