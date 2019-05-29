package ziface

type IResponse interface {
	// 响应客户端数据
	Send(id uint32, data []byte) error

	// 带有缓冲的响应
	SendBuf(id uint32, data []byte) error

	// 监听通道中写入的数据
	Process()

	// 关闭数据通道
	Close()
}
