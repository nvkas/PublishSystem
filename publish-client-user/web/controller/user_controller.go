package controller

import (
	"context"
	"flag"
	"fmt"
	"github.com/kataras/iris"
	"github.com/smallnest/rpcx/client"
	"log"
	"publish_client_user/datamodels"
	"publish_client_user/tool"
	"publish_client_user/web/middleware"
	"time"
)

var (
	consulAddr = flag.String("consulAddr", tool.ConfJson.ConsulAddr, "consul address")
	basePath   = flag.String("base", tool.ConfJson.BasePath, "prefix path")
	xClient    client.XClient
)

func init() {
	flag.Parse()
	d := client.NewConsulDiscovery(*basePath, "zzy.server.user", []string{*consulAddr}, nil)
	xClient = client.NewXClient("zzy.server.user", client.Failtry, client.RandomSelect, d, client.DefaultOption)
}
func GetClient() client.XClient {
	return xClient
}

//统一连接服务端处理方法
func Call(serviceMethod string, maps map[string]interface{}, result *datamodels.Result) *datamodels.Result {
	err := GetClient().Call(context.Background(), serviceMethod, maps, &result)
	if err != nil {
		log.Printf("ERROR failed to call: %v", err)
		return result
	}
	return result
}
func CallPage(serviceMethod string, maps map[string]interface{}, result *datamodels.ResultPage) *datamodels.ResultPage {
	err := GetClient().Call(context.Background(), serviceMethod, maps, &result)
	if err != nil {
		log.Printf("ERROR failed to call: %v", err)
		return result
	}
	return result
}

type UserController struct {
	Ctx iris.Context
}

func (c *UserController) PostLogin() (result *datamodels.Result) {
	var user datamodels.User
	c.Ctx.ReadJSON(&user)
	fmt.Println(user)
	err := GetClient().Call(context.Background(), "Login", &user, &result)
	if err != nil {
		log.Printf("ERROR failed to call: %v", err)
		return result
	} else {
		if result.Status == false {
			return result
		}
		//创建客户端对应cookie以及在服务器中进行记录
		var sessionID = middleware.SMgr.StartSession(c.Ctx.ResponseWriter(), c.Ctx.Request())
		log.Println("-------------------创建新的sessionID:", sessionID)
		//192.168.0.115
		ip := c.Ctx.RemoteAddr()
		user.RemoteAddr = ip
		user.Session = sessionID
		user.AccessTime = time.Now()
		var loginUserInfo = user
		//踢除重复登录的
		var onlineSessionIDList = middleware.SMgr.GetSessionIDList()
		for _, onlineSessionID := range onlineSessionIDList {
			log.Println("-------------------onlineSessionID:", onlineSessionID)
			if userInfo, ok := middleware.SMgr.GetSessionVal(onlineSessionID, "UserInfo"); ok {
				if value, ok := userInfo.(datamodels.User); ok {
					if value.ID == user.ID {
						fmt.Println("-------------------踢除重复登录SessionID:", onlineSessionID)
						middleware.SMgr.EndSessionBy(onlineSessionID)
					}
				}
			}
		}
		//设置变量值
		middleware.SMgr.SetSessionVal(sessionID, "UserInfo", loginUserInfo)
		tokenString := middleware.GenerateToken(&user)
		result.Token = tokenString
		return result
	}
}
