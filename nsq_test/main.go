package main

import (
	"nsq_test/tools"
	"strconv"
	"time"
)

//业务代码，生产消息
func main() {

	for i := 0; i < 5; i++ { //循环调用
		param := tools.Param{
			Addr:  "127.0.0.1:4150",
			Topic: "test",
			Msg:   []byte(strconv.Itoa(i)),
		}
		tools.InitProducer(param)
	}

	tools.InitCustomer("test","ch","127.0.0.1:4150")
	time.Sleep(time.Second * 3)

}
