package ziface

type IRequest interface {
	GetConn() IConnection

	GetMessageId() uint32

	GetData() []byte

	GetDataLen() uint32

	// 解析数据报
	ParseMessage() error

	// 开始处理客户端请求
	Process()
}
