package tool

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

//初始化配置文件
func Load(filename string, v interface{}) {
	//ReadFile函数会读取文件的全部内容，并将结果以[]byte类型返回
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return
	}
	//读取的数据为json格式，需要进行解码
	err = json.Unmarshal(data, v)
	if err != nil {
		return
	}
}

type JsonConf struct {
	MysqlName     string
	MysqlHost     string
	MysqlPort     string
	MysqlPassword string
	MysqlDatabase string
	PublishPort   string
	ConsulAddr    string
	NsqAddr       string
	BasePath      string
	Ip            string
}

var ConfJson JsonConf

func init() {
	Load("./conf.json", &ConfJson)
	fmt.Println(ConfJson)
}
