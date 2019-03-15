package tool

import (
	"fmt"
	"github.com/nsqio/go-nsq"
)

var Producer *nsq.Producer

// 初始化生产者
func InitProducer(str string) {
	var err error
	fmt.Println("address: ", str)
	Producer, err = nsq.NewProducer(str, nsq.NewConfig())
	if err != nil {
		panic(err)
	}
}

//发布消息
func Publish(topic string, message string) error {
	var err error
	if Producer != nil {
		if message == "" { //不能发布空串，否则会导致error
			return nil
		}
		err = Producer.Publish(topic, []byte(message)) // 发布消息
		return err
	}
	return fmt.Errorf("producer is nil", err)
}