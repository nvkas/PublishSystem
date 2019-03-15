package tool

var RpcConfig = &RpcConf{}

type RpcConf struct {
	Addr     string `json:"Addr"`
	ConsulAddr string `json:"ConsulAddr"`
	BasePath string `json:"BasePath"`
	NsqAddr string `json:"NsqAddr"`
	UserBasePath string `json:"UserBasePath"`
	UserServicePath string `json:"UserServicePath"`
}

func init() {
	JsonParse := NewJsonStruct()
	JsonParse.Load("./config.json", RpcConfig)
}
