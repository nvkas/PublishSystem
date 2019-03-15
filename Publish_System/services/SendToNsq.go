package services

import (
	"encoding/json"
	"errors"
	"github.com/nsqio/go-nsq"
)

func SendToNsq(data map[string]interface{}, nsqAddress string) (err error) {

	//将数据转换成json
	msg , err := json.Marshal(data)
	if err != nil {
		return
	}

	//创建生产者
	producer , err := nsq.NewProducer(nsqAddress,nsq.NewConfig())
	if err != nil {
		return
	}

	//如果msg为空,会直接panic
	if string(msg) == "" {
		return errors.New("msg cant be none")
	}

	err = producer.Publish("testtopic",msg)
	if err != nil {
		return
	}


	return nil
}
