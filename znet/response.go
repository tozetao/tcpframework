package znet

import (
	"fmt"
	"zinx/ziface"
	"zinx/zutils"
)

type Response struct {
	conn ziface.IConnection

	// 数据通道
	dataChan chan []byte

	// 带有缓冲的数据通道高
	dataBufChan chan []byte
}

func NewResponse(conn ziface.IConnection) *Response {
	r := &Response{
		conn:        conn,
		dataChan:    make(chan []byte),
		dataBufChan: make(chan []byte, zutils.ZINX.DataBufSize),
	}
	return r
}

// 发送数据给对方。其实应该有一个Response对象来负责响应的事情。
func (r *Response) Send(id uint32, data []byte) error {
	dp := NewDataPack()
	message, err := dp.Pack(NewMessage(id, data))
	if err != nil {
		return err
	}
	r.dataChan <- message
	return nil
}

func (r *Response) SendBuf(id uint32, data []byte) error {
	dp := NewDataPack()
	message, err := dp.Pack(NewMessage(id, data))
	if err != nil {
		return err
	}
	r.dataBufChan <- message
	return nil
}

// 处理客户端的响应：监听DataChan通道是否有数据，如果有就发送给客户端
func (r *Response) Process() {
loop:
	for {
		select {
		case data, ok := <-r.dataChan:
			if ok {
				r.write(data)
			} else {
				break loop
			}
		case data, ok := <-r.dataBufChan:
			if ok {
				fmt.Println("receive data from dataBufChan")
				r.write(data)
			} else {
				fmt.Println("close dataBufChan")
				break loop
			}
		}
	}
}

func (r *Response) write(data []byte) {
	if _, err := r.conn.GetTCPConn().Write(data); err != nil {
		fmt.Println("send data err: ", err)
	}
}

func (r *Response) Close() {
	fmt.Println("response data chan close.")
	close(r.dataChan)
	close(r.dataBufChan)
}
