package typemethod

import (
	"fmt"
	"publish_server_core/tool"
	"publish_server_core/datamodels"
	"github.com/smallnest/rpcx/log"
	"reflect"
	"strings"
)

type MethodMapsType map[uint]reflect.Value

var MethodMaps MethodMapsType

func init() {
	//添加接口只需添加对应的map及方法即可
	projectTypeMap := make(map[uint]string)
	projectTypeMap[0] = "ExeDefault"
	projectTypeMap[1] = "ExeVue"
	projectTypeMap[2] = "ExeGo"
	projectTypeMap[3] = "ExeGo"

	var ru Routers

	crMap := make(MethodMapsType, 0)
	//创建反射变量，注意这里需要传入ru变量的地址；
	//不传入地址就只能反射Routers静态定义的方法
	vf := reflect.ValueOf(&ru)
	vft := vf.Type()
	//读取方法数量
	mNum := vf.NumMethod()
	//fmt.Println("NumMethod:", mNum)
	//遍历路由器的方法，并将其存入控制器映射变量中
	for i := 0; i < mNum; i++ {
		mName := vft.Method(i).Name
		//fmt.Println("index:", i, " MethodName:", mName)
		for k,v := range projectTypeMap{
			if v == mName {
				crMap[k] = vf.Method(i)
			}
		}
	}
	MethodMaps = crMap
}

//定义路由器结构类型
type Routers struct {
}

//默认调用
func (this *Routers) ExeDefault(afSetting *datamodels.Setting,puSetting *datamodels.Setting,flag bool,beSetting *datamodels.Setting,project datamodels.Project,runPath string,pullPath string,projectName string) {
	log.Error("没有获取到项目类型")
}

//调用Vue
func (this *Routers) ExeVue(afSetting *datamodels.Setting,puSetting *datamodels.Setting,flag bool,beSetting *datamodels.Setting,project datamodels.Project,runPath string,pullPath string,projectName string) {
	fmt.Println("------ExeVue------")
	//添加配置文件,需要配置文件路径
	//tool.WriteFile(nginxConfPath,project.Name+".conf",tool.GetNginxTemplate(runPath,"9527",[]string{"/"+projectName}))

	//刷新配置文件
	//flag = tool.SplitCommandAndExec("cd "+nginxInstallPath+"||./nginx -s reload")

	//打包
	//beSetting.Values = "cd "+pullPath+projectName+"||npm run prod"
	beSetting.Values = "cd "+pullPath+projectName + "||npm install||npm run build:prod"
	beSetting.Values += "||mv dist/ "+runPath
	//beSetting.Values += "||cp –r dist/ "+runPath
}

//调用Go
func (this *Routers) ExeGo(afSetting *datamodels.Setting,puSetting *datamodels.Setting,flag bool,beSetting *datamodels.Setting,project datamodels.Project,runPath string,pullPath string,projectName string) {
	fmt.Println("------ExeGo------")
	//杀进程
	if flag {
		pid := tool.GetPidByProcessName(projectName)
		if pid != "" {
			flag = tool.SplitCommandAndExec("kill -9 "+pid)
		}
	}
	//复制配置文件ConfAddr
	if project.ConfAddr != "" {
		log.Info("复制配置文件")
		confAddr := strings.Split(project.ConfAddr, ",")
		for _,conf := range confAddr {
			pack := ""	//config.json	conf/config.json
			if strings.Contains(conf, "/") {
				pack = conf[ : strings.LastIndex(conf,"/")]	//conf
				exist2,_ := tool.PathExists(runPath+pack)
				if !exist2 {
					if !tool.MkDirs(runPath+pack) {
						log.Error("创建配置文件文件夹路径失败!")
					}
				}
			}

			flag = tool.SplitCommandAndExec("cp "+pullPath+projectName+"/"+conf+" "+runPath+pack)
			if !flag {
				break
			}
		}
	}

	//发布前命令:结尾添加移动二进制文件
	beSetting.Values += "cd "+pullPath+projectName+"||go build -o "+runPath+projectName	//+"||chmod +x "+runPath+projectName
	//beSetting.Values = "cd "+pullPath+projectName+"||"+beSetting.Values
	//是否存在二进制文件
	//exist,_ = tool.PathExists(pullPath+projectName+"/"+projectName)
	//log.Info(pullPath+projectName+"/"+projectName)
	//if exist {
	//	beSetting.Values += "||chmod +x "+projectName
	//	if projectName != "" && project.User.ProjectWorkPath != "" {
	//		beSetting.Values += "||mv -f "+pullPath+projectName+"/"+projectName+" "+runPath
	//	}
	//}

	//发布命令:头部添加cd到目录
	puSetting.Values += "cd "+runPath+"||nohup ./"+projectName+" &||"+puSetting.Values

}