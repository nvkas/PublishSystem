package service

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/kataras/iris/core/errors"
	"github.com/nsqio/go-nsq"
	"github.com/smallnest/rpcx/log"
	"publish_server_core/datamodels"
	"publish_server_core/repositories"
	"publish_server_core/services/typemethod"
	"publish_server_core/tool"
	"reflect"
	"strings"
)

var repoSetting = repositories.NewSettingRepository()

type SettingServer struct {
}

func init() {
	InitConsumer("publish_core_start", "publish_core_user001", tool.RpcConfig.NsqAddr,"start")
	InitConsumer("publish_core_stop", "publish_core_user002", tool.RpcConfig.NsqAddr,"stop")
}

// 消费者
type ConsumerStart struct{}
type ConsumerStop struct{}

//处理Hook消息
func (*ConsumerStart) HandleMessage(msg *nsq.Message) error {
	fmt.Println("INFO:Body:",string(msg.Body))
	var maps map[string]interface{}
	json.Unmarshal(msg.Body,&maps)

	if err := Start(maps);err != nil {
		log.Error(err)
	}
	return nil
}
//处理Hook消息
func (*ConsumerStop) HandleMessage(msg *nsq.Message) error {
	fmt.Println("INFO:Body:",string(msg.Body))
	var maps map[string]interface{}
	json.Unmarshal(msg.Body,&maps)

	if err := Stop(maps);err != nil {
		log.Error(err)
	}
	return nil
}

//初始化消费者
func InitConsumer(topic string, channel string, address string, method string) {
	fmt.Println("----------Init"+method+"Consumer----------------")
	cfg := nsq.NewConfig()
	//cfg.LookupdPollInterval = time.Second          //设置重连时间
	c, err := nsq.NewConsumer(topic, channel, cfg) // 新建一个消费者
	if err != nil {
		log.Error(err)
	}
	c.SetLogger(nil, 0)        //屏蔽系统日志
	if method == "stop" {
		c.AddHandler(&ConsumerStop{}) // 添加消费者接口
	}else if method == "start"{
		c.AddHandler(&ConsumerStart{}) // 添加消费者接口
	}

	//建立NSQLookupd连接
	if err := c.ConnectToNSQLookupd(address); err != nil {
		log.Error(err)
	}
}

func (t *SettingServer) ExecAllCommand(ctx context.Context, maps *map[string]interface{}, result *datamodels.Result) error {
	fmt.Println("-----------ExecAllCommand-----------")
	name := fmt.Sprintf("%v",(*maps)["repository_name"])
	log.Info("repository_name:"+name)
	if name == "" {
		return errors.New("执行失败!!")
	}
	project,systemSettings := repoSetting.GetProjectByName(name)
	if project.WarehouseName == "" || project.GitAddress == "" {
		return errors.New("执行失败!")
	}
	var beSetting datamodels.Setting
	var puSetting datamodels.Setting
	var afSetting datamodels.Setting

	settings := project.Setting
	for _,sett := range settings{
		if sett.Keys == "beforeCommand" {
			beSetting = sett
		}else if sett.Keys == "publishCommand" {
			puSetting = sett
		}else if sett.Keys == "afterCommand" {
			afSetting = sett
		}
	}
	goPath := ""
	vuePath := ""
	projectPath := ""
	for _,systemSetting := range systemSettings{
		if systemSetting.Keys == "GoPath" {
			goPath = systemSetting.Values
		}else if systemSetting.Keys == "VuePath" {
			vuePath = systemSetting.Values
		}else if systemSetting.Keys == "ProjectPath" {
			projectPath = systemSetting.Values
		}
	}

	projectName := project.WarehouseName[strings.Index(project.WarehouseName, "/")+1 :]
	pullPath := goPath
	//前端项目
	if project.ProjectTypeId == 1 {
		pullPath = vuePath
	}
	if !strings.HasSuffix(pullPath, "/"){
		pullPath += "/"
	}
	if !strings.HasSuffix(projectPath, "/"){
		projectPath += "/"
	}
	runPath := projectPath+project.User.LoginName+"/"+project.Name+"/"
	exist,_ := tool.PathExists(runPath)
	if !exist {
		if !tool.MkDirs(runPath) {
			log.Error("创建二进制执行文件文件夹路径失败!")
		}
	}

	//pullPath += project.User.LoginName+"/"+"src"+"/"
	//判断是否已经CLone下来了
	exist,_ = tool.PathExists(pullPath+projectName)

	//是否已存在
	flag := true
	if !exist {
		//Clone代码
		fmt.Println("-----项目clone目录------")
		fmt.Println(pullPath)
		fmt.Println("-----------")
		flag = tool.CloneRepo(pullPath,"master",project.GitAddress)
	}

	//pull代码
	text := "cd "+pullPath+projectName+" ||"+"git pull"
	if flag {
		flag = tool.SplitCommandAndExec(text)
	}

	//afSetting,puSetting,flag ,beSetting ,nginxConfPath ,project ,runPath ,nginxInstallPath ,pullPath ,projectName
	parms := []reflect.Value{reflect.ValueOf(&afSetting),reflect.ValueOf(&puSetting),reflect.ValueOf(flag),reflect.ValueOf(&beSetting),
		reflect.ValueOf(project),reflect.ValueOf(runPath),reflect.ValueOf(pullPath),reflect.ValueOf(projectName)}
	//使用类型调用指定方法
	typemethod.MethodMaps[project.ProjectTypeId].Call(parms)

	////前端项目
	//if project.ProjectTypeId == 1 {
	//
	////后端项目
	//}else {
	//
	//}

	if flag {
		err := t.ExecBeforeCommand(ctx,&beSetting,result)
		if err == nil {
			err = t.ExecPublishCommand(ctx,&puSetting,result)
			if err == nil {
				err = t.ExecAfterCommand(ctx,&afSetting,result)
				if err == nil {
					result.Data = "success"
					result.Status = true
					result.Code = "200"
					return nil
				}
				fmt.Println(err)
				return nil
			}
			fmt.Println(err)
		}
		fmt.Println(err)
	}

 	return errors.New("执行失败")
}
func (t *SettingServer) ExecBeforeCommand(ctx context.Context, setting *datamodels.Setting, result *datamodels.Result) error {
	fmt.Println("-----------ExecBeforeCommand-----------")
	//登陆Git(可一次操作)
	//Pull代码

	//打包
	//剪切到运行路径覆盖原打包文件
	//杀进程

	//按照换行符分割
	//text := "cd /home/go/src/qiaoyi_back/ \r\n" +
	//"git pull \r\n"+
	//"sudo fuser -k -n tcp 8080 \r\n"+
	//"go build main.go \r\n"+
	//"chmod +x main"+
	//"mv -f main /home/go/src/test/"
	text := setting.Values
	flag := tool.SplitCommandAndExec(text)
	if flag {
		result.Data = "success"
		result.Status = true
		result.Code = "200"
		return nil
	}
	return errors.New("执行发布前命令失败")
}
func (t *SettingServer) ExecPublishCommand(ctx context.Context, setting *datamodels.Setting, result *datamodels.Result) error {
	fmt.Println("-----------ExecPublishCommand-----------")
	//nohup ./main &
	//text := "cd /home/go/src/test/ \r\n"+
	//	"nohup ./main"
	text := setting.Values
	flag := tool.SplitCommandAndExec(text)
	if flag {
		result.Data = "success"
		result.Status = true
		result.Code = "200"
		return nil
	}
	return errors.New("执行发布命令失败")
}
func (t *SettingServer) ExecAfterCommand(ctx context.Context, setting *datamodels.Setting, result *datamodels.Result) error {
	fmt.Println("-----------ExecAfterCommand-----------")
	text := setting.Values
	flag := tool.SplitCommandAndExec(text)
	if flag {
		result.Data = "success"
		result.Status = true
		result.Code = "200"
		return nil
	}
	return errors.New("执行发布后命令失败")
}
//运行项目
func Start(maps map[string]interface{}) error {
	projectName := fmt.Sprintf("%v",maps["warehouseName"])
	if strings.Contains(projectName, "/"){
		projectName = projectName[strings.Index(projectName, "/")+1 :]
	}
	projectPath := fmt.Sprintf("%v",maps["projectAddr"])
	if !strings.HasSuffix(projectPath, "/"){
		projectPath += "/"
	}
	userName := fmt.Sprintf("%v",maps["userName"])

	text := "cd " + projectPath + userName + "/" + projectName +"||nohup ./"+projectName + " &"
	flag := tool.SplitCommandAndExec(text)
	if flag {
		return nil
	}

	return errors.New("程序运行失败")
}
//终止项目
func Stop(maps map[string]interface{}) error {
	projectName := fmt.Sprintf("%v",maps["warehouseName"])
	if strings.Contains(projectName, "/"){
		projectName = projectName[strings.Index(projectName, "/")+1 :]
	}
	pid := tool.GetPidByProcessName(projectName)
	if pid != "" {
		flag := tool.SplitCommandAndExec("kill -9 "+pid)
		if flag {
			return nil
		}
	}
	return errors.New("程序终止失败")
}