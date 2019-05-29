package znet

import (
	"fmt"
	"io"
	"zinx/ziface"
	"zinx/zutils"
)

type Request struct {
	conn    ziface.IConnection
	message ziface.IMessage
}

func NewRequest(conn ziface.IConnection) *Request {
	return &Request{
		conn:    conn,
		message: nil,
	}
}

func (r *Request) GetConn() ziface.IConnection {
	return r.conn
}

func (r *Request) GetData() []byte {
	return r.message.GetData()
}

func (r *Request) GetDataLen() uint32 {
	return r.message.GetDataLen()
}

func (r *Request) GetMessageId() uint32 {
	return r.message.GetId()
}

// 定义一个处理TCP连接的方法：读取TCP连接数据，并封装到Message对象中。
func (r *Request) ParseMessage() error {
	dp := NewDataPack()
	tcpConnection := r.GetConn().GetTCPConn()

	// 解包头
	headBuf := make([]byte, dp.GetHeadLen())
	if _, err := io.ReadFull(tcpConnection, headBuf); err != nil {
		return err
	}

	message, err := dp.UnPack(headBuf)
	if err != nil {
		return err
	}

	// 读取包体数据
	if message.GetDataLen() > 0 {
		data := make([]byte, message.GetDataLen())
		if _, err := io.ReadFull(tcpConnection, data); err != nil {
			return err
		}
		message.SetData(data)
	}

	// 将解包的Message赋值给Request对象
	r.message = message
	return nil
}

// 读取连接数据，并交由对应的路由处理
func (r *Request) Process() {
	defer func() {
		r.conn.Stop()
	}()

	conn := r.conn

	for {
		err := r.ParseMessage()

		// 断开连接
		if err == io.EOF {
			break
		}
		// 报错
		if err != nil {
			fmt.Println("request processing err: ", err)
			break
		}

		if zutils.ZINX.WorkerPoolSize > 0 {
			conn.GetMessageHanlder().AddTask(conn)
		} else {
			conn.GetMessageHanlder().Dispatch(r, conn.GetResponse())
		}
	}
}
