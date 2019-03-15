package tool

import (
	//"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"Publish_System/services"
)

var config *Config

func NewHttpServer(conf *Config)  {
	config = conf
	http.HandleFunc("/PushSignal",ReceiveHook)

	err := http.ListenAndServe(config.HttpAddress, nil)
	if err != nil {
		log.Fatalf("HttpServer init failed:%s",err)
	}

}

func ReceiveHook(w http.ResponseWriter, r *http.Request)  {
	//读取数据
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatalf("read data failed:%s",err)
	}


	//处理数据
	new_data := services.DealWithData(data)

	//发送信号到nsq
	//err = services.SendToNsq(new_data,config.NsqAddress)
	//if err != nil {
	//	如果出错,发送邮件通知
		//services.SendEmail(fmt.Sprintf("有用户推送新代码,但是发到Nsq失败:%s",err),config.FormAddress,config.ReceiveAddress,config.UserName,config.AuthorizationCode)
		//如果用fatalf会导致服务直接退出,此处不希望退出
		//log.Printf("Send msg failed:%s",err)
		//return
	//}

	//发送信号
	err = services.SendSignal("server.core.allCommand","/core",config.ConsulAddress,new_data)
	if err != nil {
		log.Printf("SendSingnal failed:%s",err)
	}

	return


}
