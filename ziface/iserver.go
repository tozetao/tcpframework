package ziface

type IServer interface {
	// 启动服务器
	Start()

	// 关闭服务器
	Stop()

	// 开启业务服务方法
	Serve()

	// 添加路由
	AddRouter(messageId uint32, router IRouter)

	GetConnManager() IConnManager

	// 设置连接启动时的回调函数
	SetOnConnStart(hock func(IConnection))

	// 设置断开连接时的回调函数
	SetOnConnStop(hock func(IConnection))

	CallOnConnStart(conn IConnection)

	CallOnConnStop(conn IConnection)
}
