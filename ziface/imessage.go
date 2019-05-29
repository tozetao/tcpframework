package ziface

type IMessage interface {
	// 返回消息Id
	GetId() uint32

	// 获取数据长度
	GetDataLen() uint32

	// 获取具体数据
	GetData() []byte

	SetId(id uint32)

	SetDataLen(len uint32)

	SetData(data []byte)
}
