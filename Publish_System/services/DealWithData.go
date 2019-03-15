package services

import (
	"encoding/json"
	"log"
)

func DealWithData(data []byte) map[string]interface{} {
	//将数据转换成map做处理
	temp := make(map[string]interface{})

	//即将返回的map
	new_data := make(map[string]interface{})
	err := json.Unmarshal(data,&temp)
	if err != nil {
		log.Fatalf("data invalid:%s",err)
	}

	//断言转换
	temp_map , _ := temp["repository"].(map[string]interface{})
	new_data["repository_name"] = temp_map["full_name"]

	temp_map , _ = temp["pusher"].(map[string]interface{})
	new_data["pusher"] = temp_map["username"]

	return new_data
}
