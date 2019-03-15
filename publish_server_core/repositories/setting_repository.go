package repositories

import (
	"context"
	"flag"
	"fmt"
	"github.com/goinggo/mapstructure"
	"github.com/smallnest/rpcx/client"
	"publish_server_core/datamodels"
	"publish_server_core/tool"
)

type SettingRepository interface {
	GetProjectByName(name string) (project datamodels.Project,systemSetting []datamodels.SystemSetting)
	GetProjectTypeAll() (projectTypes []datamodels.ProjectType)
}

func NewSettingRepository() SettingRepository {
	return &settingRepository{}
}

type settingRepository struct {}

var userClient client.XClient

func init() {
	flag.Parse()
	d := client.NewConsulDiscovery(tool.RpcConfig.UserBasePath, tool.RpcConfig.UserServicePath, []string{tool.RpcConfig.ConsulAddr}, nil)
	userClient = client.NewXClient(tool.RpcConfig.UserServicePath, client.Failtry, client.RandomSelect, d, client.DefaultOption)
}

func GetUserClient() client.XClient {
	return userClient
}

func (m *settingRepository) GetProjectByName(name string) (project datamodels.Project,systemSettings []datamodels.SystemSetting) {
	fmt.Println("-------GetProjectByName---------")

	maps := make(map[string]interface{})
	maps["repository_name"] = name

	var result datamodels.Result

	err := GetUserClient().Call(context.Background(), "FindWarehouseProject", &maps, &result)
	if err != nil {
		fmt.Println("failed to call: %v", err)
	}
	fmt.Println("------返回值-------")
	data := result.Data.(map[string]interface{})
	err = mapstructure.Decode(data["project"], &project)
	err = mapstructure.Decode(data["systemSetting"], &systemSettings)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(project)
	//defer GetUserClient().Close()

	return
}

func (m *settingRepository) GetProjectTypeAll() (projectTypes []datamodels.ProjectType){
	fmt.Println("-------GetProjectTypeAll---------")
	maps := make(map[string]interface{})
	var result datamodels.Result

	err := GetUserClient().Call(context.Background(), "FindProjectTypeAll", &maps, &result)
	if err != nil {
		fmt.Println("failed to call: %v", err)
		//log.Error(errors.New(fmt.Sprintf("failed to call: %v", err)))
	}
	err = mapstructure.Decode(result.Data, &projectTypes)
	if err != nil {
		fmt.Println(err)
	}
	return
}