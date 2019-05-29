package znet

import (
	"net"
	"zinx/ziface"
)

type Connection struct {
	Conn *net.TCPConn

	ConnId int32

	isClosed bool

	MessageHandler ziface.IMessageHandler

	// Request对象
	request ziface.IRequest

	// response对象
	response ziface.IResponse

	server ziface.IServer
}

//创建连接的方法
func NewConntion(s ziface.IServer, conn *net.TCPConn, id int32, msgHandler ziface.IMessageHandler) *Connection {
	c := &Connection{
		Conn:           conn,
		ConnId:         id,
		isClosed:       false,
		MessageHandler: msgHandler,
		request:        nil,
		response:       nil,
		server:         s,
	}

	c.request = NewRequest(c)
	c.response = NewResponse(c)
	s.GetConnManager().Add(id, c)

	return c
}

// 处理客户端连接
func (c *Connection) Start() {
	go c.request.Process()
	go c.response.Process()
	c.server.CallOnConnStart(c)
}

// 关闭链接
func (c *Connection) Stop() {
	c.server.CallOnConnStop(c)

	c.isClosed = true

	// 关闭TCP连接
	c.Conn.Close()

	// 关闭Response对象
	c.response.Close()

	// 移除连接
	c.server.GetConnManager().Remove(c.ConnId)
}

func (c *Connection) GetConnId() int32 {
	return c.ConnId
}

func (c *Connection) GetTCPConn() *net.TCPConn {
	return c.Conn
}

func (c *Connection) RemoteAddr() net.Addr {
	return c.Conn.RemoteAddr()
}

func (c *Connection) GetRequest() ziface.IRequest {
	return c.request
}

func (c *Connection) GetResponse() ziface.IResponse {
	return c.response
}

func (c *Connection) GetMessageHanlder() ziface.IMessageHandler {
	return c.MessageHandler
}
