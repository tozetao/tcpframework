package znet

import (
	"fmt"
	"io"
	"net"
	"testing"
	"time"
	"zinx/ziface"
)

type PingRouter struct {
	BaseRouter
}

func (this *PingRouter) Handle(request ziface.IRequest, response ziface.IResponse) {
	fmt.Println("Call Router Handle()")
	fmt.Printf("read from %d bytes, data: %s\n", request.GetDataLen(), request.GetData())

	err := response.Send(1, []byte("ping...ping"))
	if err != nil {
		fmt.Println("Call Handle error", err)
	}
}

func ClientTest() {
	// 保证服务器先运行
	time.Sleep(3 * time.Second)
	fmt.Println("Client Test... start")

	conn, err := net.Dial("tcp", "127.0.0.1:9999")
	defer conn.Close()
	if err != nil {
		fmt.Println("client start err, exit!")
		return
	}

	//发送一次数据
	dp := NewDataPack()
	msg, err := dp.Pack(NewMessage(1001, []byte("hello world")))
	if err != nil {
		fmt.Println("client pack err: ", err)
		return
	}
	if _, err := conn.Write(msg); err != nil {
		fmt.Println("client write err: ", err)
		return
	}

	//读取一次数据
	headBuf := make([]byte, dp.GetHeadLen())
	if _, err := io.ReadFull(conn, headBuf); err != nil {
		fmt.Println("client headBuf err: ", err)
		return
	}

	message, err := dp.UnPack(headBuf)
	if err != nil {
		fmt.Println("client UnPack headBuf err: ", err)
		return
	}

	// 读取包体数据
	if message.GetDataLen() > 0 {
		data := make([]byte, message.GetDataLen())
		if _, err := io.ReadFull(conn, data); err != nil {
			fmt.Println("client UnPack data err: ", err)
			return
		}
	}
	fmt.Printf("read %d bytes from server, data: %s\n", message.GetDataLen(), message.GetData())

	time.Sleep(3 * time.Second)
}

func Send(in chan []byte) {
	bytes := []byte("hello world")
	in <- bytes
}

func Receive(out chan []byte) {
loop:
	for {
		select {
		case data, ok := <-out:
			if !ok {
				break loop
			}
			fmt.Println(data)
		}
	}
	fmt.Println("end")
}

func TestServer(t *testing.T) {
	// go ClientTest()
	// s := NewServer("zinx v1.0")
	// s.AddRouter(0, &PingRouter{})
	// s.Serve()

	dataChan := make(chan bool)

	go func() {
		time.Sleep(3 * time.Second)
		dataChan <- true
	}()

	go func() {
		close(dataChan)
	}()

	time.Sleep(10 * time.Second)
}

/*
Server
	Server负责服务管理。

Connection
	封装处理TCP连接的业务逻辑。

Request
	封装Connection对象和请求数据

Router
	路由是映射客户端操作的对象，它定义了如何处理客户端的某种操作，由框架使用者来负责实现路由。

Message
	封装处理数据报的逻辑。
	len（4个字节） + id（4个字节） + data（由len字段指定的长度）

DataPack
	提供拆包与解包

MessageManager
	消息管理对象，实现路由注册、调用路由来处理消息。


Server:AddRouter -> Connection -> 调用路由对象，执行具体的逻辑处理
*/

/*

在Connection对象中开启俩个协程：
	读协程：
		除了负责读取客户端数据，还负责把数据交给对应的路由处理。

	写协程：
		从通道中读取数据，然后写入到客户端连接中。

	好处：首先读写分离会使逻辑更加清晰，另外写入数据由另外一个协程负责，处理客户端次数的吞吐量会提高。


	exitChan bool：该通道负责负责通知主协程，客户端连接已经关闭。
	DataChan []byte: 负责从Read协程向Write协程传递数据的通道

	Read -> Write -> Client

	Write
		在写协程中，会使用Select监听俩个通道，一个exitChan、一个dataChan。
		如果exitChan有数据，代表着需要关闭客户端连接，以及相关结束相关协程的运行；
		如果使dataChan有数据，代表着需要写入数据给客户端。

	Read协程
		如果客户端发送的数据报异常或者断开连接，那么会像exitChan发送一个布尔值

*/

/*
消息队列与Worker池
	Worker池是若干数量协程的集合，每个Worker协程对应一个任务。
	在本框架中的任务是指业务逻辑的处理，即worker负责处理业务逻辑。


	消息队列
		存储着待处理的Message。

	Connection会将Message写入到消息队列中，由对应的Worker去处理。

	消息队列的大小限制了同一个时间能够处理请求的业务逻辑数量，好处是避免了Connection协程的上下文切换。处理
	业务逻辑的切换转移到Worker池了。
*/
