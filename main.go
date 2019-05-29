package main

import (
	"fmt"
	"time"
	"zinx/ziface"
	"zinx/znet"
)

type PingRouter struct {
	znet.BaseRouter
}

// 类型0的消息
func (this *PingRouter) Handle(request ziface.IRequest, response ziface.IResponse) {
	fmt.Println("receive message id: ", request.GetMessageId())

	time.Sleep(5 * time.Second)

	err := response.SendBuf(0, []byte("ping...ping..."))
	if err != nil {
		fmt.Println("Call PingRouter error", err)
	}
}

// 类型1的消息
type HelloRouter struct {
	znet.BaseRouter
}

func (this *HelloRouter) Handle(request ziface.IRequest, response ziface.IResponse) {
	fmt.Println("receive message id: ", request.GetMessageId())
	err := response.SendBuf(1, []byte("hello"))
	if err != nil {
		fmt.Println("Call HelloRouter error", err)
	}
}

func DoConnBegin(conn ziface.IConnection) {
	var content string
	content = string(conn.GetConnId()) + " do conn begin"
	conn.GetResponse().SendBuf(1, []byte(content))
}

func DoConnEnd(conn ziface.IConnection) {
	fmt.Println(conn.GetConnId(), " do conn end.")
}

func main() {
	s := znet.NewServer("zinx v1.0")

	// 设置钩子
	// s.SetOnConnStart(DoConnBegin)
	// s.SetOnConnStop(DoConnEnd)

	// 添加路由
	s.AddRouter(0, &PingRouter{})
	s.AddRouter(1, &HelloRouter{})
	s.Serve()
}
