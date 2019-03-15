package tool

import (
	"github.com/rcrowley/go-metrics"
	"github.com/smallnest/rpcx/server"
	"github.com/smallnest/rpcx/serverplugin"
	"log"
	"time"
)

func NewServer(ServerAddress string,ConsulAddress []string,BasePath string) *server.Server {
	s := server.NewServer()
	addRegistryPlugin(s, ServerAddress,ConsulAddress,BasePath)
	return s
}

func addRegistryPlugin(s *server.Server, ServerAddress string,ConsulAddress []string,BasePath string) {
	r := &serverplugin.ConsulRegisterPlugin{
		ServiceAddress: "tcp@" + ServerAddress,
		ConsulServers:  ConsulAddress,
		BasePath:       BasePath,
		Metrics:        metrics.NewRegistry(),
		UpdateInterval: time.Minute,
	}
	err := r.Start()
	if err != nil {
		log.Fatal(err)
	}
	s.Plugins.Add(r)
}
