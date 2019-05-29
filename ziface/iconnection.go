package ziface

import "net"

/**
定义处理客户端连接的接口
*/
type IConnection interface {
	// 开始处理连接
	Start()

	// 结束连接的处理
	Stop()

	RemoteAddr() net.Addr

	GetConnId() int32

	GetTCPConn() *net.TCPConn

	GetRequest() IRequest

	GetResponse() IResponse

	GetMessageHanlder() IMessageHandler
}
