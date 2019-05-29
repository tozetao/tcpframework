package znet

import (
	"bytes"
	"encoding/binary"
	"errors"
	"zinx/ziface"
	"zinx/zutils"
)

type DataPack struct{}

func NewDataPack() *DataPack {
	return &DataPack{}
}

/*
	头部由4个字节的len和4个字节的消息id组成，因此是8个字节。
	头部的len字段指明了包体长度，id字段则是包体id，即整个消息的id。
*/
func (this *DataPack) GetHeadLen() int16 {
	return 8
}

func (this *DataPack) Pack(message ziface.IMessage) ([]byte, error) {
	buf := bytes.NewBuffer([]byte{})

	if err := binary.Write(buf, binary.LittleEndian, message.GetDataLen()); err != nil {
		return nil, err
	}

	if err := binary.Write(buf, binary.LittleEndian, message.GetId()); err != nil {
		return nil, err
	}

	if err := binary.Write(buf, binary.LittleEndian, message.GetData()); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

/*
	只是针对头部进行拆包，将返回不包括具体数据的Message对象。
*/
func (this *DataPack) UnPack(head []byte) (ziface.IMessage, error) {
	var len, id uint32
	message := &Message{}

	reader := bytes.NewReader(head)

	if err := binary.Read(reader, binary.LittleEndian, &len); err != nil {
		return nil, err
	}

	if err := binary.Read(reader, binary.LittleEndian, &id); err != nil {
		return nil, err
	}

	// 验证包体长度
	if zutils.ZINX.MaxPacketSize > 0 && len > zutils.ZINX.MaxPacketSize {
		return nil, errors.New("The body exceeds limit.")
	}

	message.SetId(id)
	message.SetDataLen(len)
	return message, nil
}
