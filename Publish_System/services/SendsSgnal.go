package services

import (
	"context"
	"github.com/smallnest/rpcx/client"
	"Publish_System/datamodels"
)

func SendSignal(serverName , basePath string,consulAddress []string,data map[string]interface{}) error  {
	cli := NewClient(serverName,basePath,consulAddress)
	err := cli.Call(context.Background(),"ExecAllCommand",data,new(datamodels.Result))
	//ExecAllCommand
	return err
}



func NewClient(Servername,BashPath string, ConsulAddress []string) client.XClient {
	consul := client.NewConsulDiscovery(BashPath, Servername, ConsulAddress, nil)
	xclient := client.NewXClient(Servername, client.Failtry, client.RandomSelect, consul, client.DefaultOption)

	return xclient

}
