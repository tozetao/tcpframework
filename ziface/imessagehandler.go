package ziface

type IMessageHandler interface {
	AddRouter(messageId uint32, router IRouter)

	// 将请求分发给对应的路由去处理
	Dispatch(request IRequest, response IResponse)

	StartWorkerPool()

	AddTask(conn IConnection)
}
