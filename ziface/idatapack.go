package ziface

/*
	定义封装、拆包相关接口
*/

type IDataPack interface {
	// 返回头部占用的字节长度
	GetHeadLen() uint16

	// 封包
	Pack(message IMessage) ([]byte, error)

	// 拆包
	UnPack(data []byte) (IMessage, error)
}
