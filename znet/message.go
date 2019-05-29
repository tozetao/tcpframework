package znet

/*
应用层协议：
	每次发送与接收的数据都是一个消息，消息头由8个字节组成，前4个字节是内容体长度，后4个字节是消息id。
*/
type Message struct {
	id   uint32
	len  uint32
	data []byte
}

func NewMessage(id uint32, data []byte) *Message {
	len := uint32(len(data))
	return &Message{
		id:   id,
		len:  len,
		data: data,
	}
}

func (m *Message) GetId() uint32 {
	return m.id
}

func (m *Message) GetDataLen() uint32 {
	return m.len
}

func (m *Message) GetData() []byte {
	return m.data
}

func (m *Message) SetId(id uint32) {
	m.id = id
}

func (m *Message) SetDataLen(len uint32) {
	m.len = len
}

func (m *Message) SetData(data []byte) {
	m.data = data
}
