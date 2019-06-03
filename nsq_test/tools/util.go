package tools

import (
	"bytes"
	"errors"
	"fmt"
	"strings"
)

//定义接口，内涵参数检查方法
type CheckOfParam interface {
	Check() error//对参数的校验函数
}

//产生消息的入参结构体
type Param struct {
	Addr  string
	Topic string
	Msg   []byte
}

//产生消息入参结构体实现参数检查接口
// 辅助函数
func (p *Param)Check()error {
	var ErrMsgs []string
	var buf bytes.Buffer
	if p.Addr=="" {
		buf.WriteString("addr can not empty....")
		ErrMsgs= append(ErrMsgs, "addr can not empty....")
	}

	if p.Topic=="" {
		buf.WriteString("addr can not empty....")
		ErrMsgs= append(ErrMsgs, "Topic can not empty....")
	}
	if p.Msg == nil{
		buf.WriteString("addr can not empty....")
		ErrMsgs= append(ErrMsgs, "Msg can not empty....")
	}
	
	//处理错误信息返回
	if len(ErrMsgs)>0 {
		errMsg:=strings.Join(ErrMsgs,"")
		return errors.New(fmt.Sprintf("Initalzing Producer Faled : %d\n",errMsg))
	}

	return nil

}
