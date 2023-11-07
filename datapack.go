package y

import (
	"bytes"
	"encoding/binary"
)

type DataPack struct {
}

var DefaultDataPack = &DataPack{}

func (d *DataPack) Pack(msg *Msg) ([]byte, error) {
	buffer := bytes.NewBuffer([]byte{})

	binary.Write(buffer, binary.BigEndian, msg.DataLen)
	binary.Write(buffer, binary.BigEndian, msg.Typ)
	binary.Write(buffer, binary.BigEndian, msg.Id)
	binary.Write(buffer, binary.BigEndian, msg.Data)

	return buffer.Bytes(), nil
}

func (d *DataPack) UnPack(data []byte) (*Msg, error) {
	buffer := bytes.NewReader(data)

	msg := &Msg{}
	err := binary.Read(buffer, binary.BigEndian, &msg.DataLen)
	if err != nil {
		return nil, err
	}
	err = binary.Read(buffer, binary.BigEndian, &msg.Typ)
	if err != nil {
		return nil, err
	}
	err = binary.Read(buffer, binary.BigEndian, &msg.Id)
	if err != nil {
		return nil, err
	}

	return msg, nil
}

func Pack(msg *Msg) []byte {
	data, _ := DefaultDataPack.Pack(msg)
	return data
}

func UnPack(data []byte) *Msg {
	msg, _ := DefaultDataPack.UnPack(data)
	return msg
}
