package tool

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

type Config struct {
	ServerName string       `json:"server_name"`                    //服务名称
	ServerAddress string    `json:"server_address"`                //服务端口
	HttpAddress string      `json:"http_address"`                 //Http服务端口
	BasePath string 	    `json:"base_path"`					 //命名空间
	NsqAddress string       `json:"nsq_address"`                //Nsq地址
	FormAddress string      `json:"form_address"`              //发件人地址
	ReceiveAddress string   `json:"receive_address"`          //收件人地址
	UserName string         `json:"user_name"`               //发送人账户
	AuthorizationCode string`json:"authorization_code"`     //发件人授权码
	ConsulAddress []string  `json:"consul_address"`        //Consul地址
	
}


func LoadConfig() *Config {
	Conf := new(Config)

	data , err := ioutil.ReadFile("./application.json")
	if err != nil {
		log.Fatalf("init config failed:%s",err)
	}

	//读取配置文件
	err = json.Unmarshal(data,Conf)
	if err != nil {
		log.Fatalf("init config failed:%s",err)
	}

	return Conf

}
