package ziface

type IConnManager interface {
	// 添加一个连接对象
	Add(connId int32, conn IConnection)

	// 移除一个连接对象
	Remove(connId int32)

	Get(connId int32) (IConnection, error)

	// 清空所有连接，并关闭所有链接。
	ClearAll()
}
