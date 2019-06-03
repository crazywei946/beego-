package tools

import (
	"errors"
	"fmt"
	"github.com/go-nsq"
)

//定义消费者结构体，需要实现nsq的接口
type Concustomer struct {

}



//从消息队列中取出 消息
//消息的处理

func InitCustomer(topic,channel,addr string) error {

	if topic=="" || channel =="" || addr=="" { //入参校验
		return errors.New("param is not current...")
	}

	//生成消费者
	con,err:=nsq.NewConsumer(topic,channel,nsq.NewConfig())
	if err != nil {
		return errors.New("Creat consumer if fale...")
	}

	//进行消息的处理函数
	con.AddHandler(&Concustomer{})


	//消费者建立nsqlookupd连接
	if err:=con.ConnectToNSQD(addr);err!=nil{
		return errors.New("con connected to nsqlookupd error...")
	}

	return nil
}

//真正的消息处理函数 ，需要实现接口
//业务处理流程均卸载此函数当中
func (con *Concustomer)HandleMessage(msg *nsq.Message) error {
	fmt.Printf("正在进行消息的处理,处理消息内容为%v\n",string(msg.Body))
	return nil
}

