package znet

import (
	"fmt"
	"strconv"
	"zinx/ziface"
	"zinx/zutils"
)

type MessageHandler struct {
	routers map[uint32]ziface.IRouter

	// 消息队列（任务队列）
	taskQunue []chan ziface.IConnection

	// Worker池的大小
	workerPoolSize uint32
}

func NewMessageHandler() *MessageHandler {
	return &MessageHandler{
		routers:        make(map[uint32]ziface.IRouter),
		taskQunue:      make([]chan ziface.IConnection, zutils.ZINX.WorkerPoolSize),
		workerPoolSize: zutils.ZINX.WorkerPoolSize,
	}
}

// 添加消息对应的路由
func (mh *MessageHandler) AddRouter(messageId uint32, router ziface.IRouter) {
	if _, ok := mh.routers[messageId]; ok {
		panic("repeated router, messageId = " + strconv.Itoa(int(messageId)))
	}
	mh.routers[messageId] = router
}

// 分发路由去处理请求
func (mh *MessageHandler) Dispatch(request ziface.IRequest, response ziface.IResponse) {
	router, ok := mh.routers[request.GetMessageId()]
	if !ok {
		fmt.Println("router of messageId is not found.")
		return
	}
	router.PreHandle(request, response)
	router.Handle(request, response)
	router.PostHandle(request, response)
}

// 启动Worker协程池
func (mh *MessageHandler) StartWorkerPool() {
	for i := 0; i < int(zutils.ZINX.WorkerPoolSize); i++ {
		mh.taskQunue[i] = make(chan ziface.IConnection)
		go mh.StartOneWorker(int32(i))
	}
}

// 启动一个Worker
func (mh *MessageHandler) StartOneWorker(workerId int32) {
	for {
		select {
		case conn := <-mh.taskQunue[workerId]:
			mh.Dispatch(conn.GetRequest(), conn.GetResponse())
		}
	}
}

// 添加任务到消息队列中
// 根据请求ID，轮询分配给Worker池处理
func (mh *MessageHandler) AddTask(conn ziface.IConnection) {
	// 计算该消息对应的WorkerId
	workerId := conn.GetConnId() % int32(zutils.ZINX.WorkerPoolSize)

	// 写入任务队列
	mh.taskQunue[workerId] <- conn
}
