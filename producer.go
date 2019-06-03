package tools

import (
	"fmt"
	"github.com/go-nsq"
)

//将消息写入到队列中
//初始化一个producer,并且产生消息
func InitProducer(p Param) error {
	//参数校验
	if err := p.Check(); err != nil {
		return err
	}
	//产生生产者
	produce, err := nsq.NewProducer(p.Addr, nsq.NewConfig())
	if err != nil {
		return err
	}
	//推送消息到队列中
	if err := produce.Publish(p.Topic, p.Msg); err != nil {
		return err
	}
	fmt.Println("消息推送完成")
	return nil

}
