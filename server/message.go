package server

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"github.com/duke-git/lancet/v2/convertor"
	"net"
)

const (
	ERROR   = 1000
	SUCCESS = 2000
	FAILURE = 5000
)

// 消息包type类型
const (
	REQUEST = iota
	RESPONSE
	NOTIFY
)

// 消息包
type MessagePacket struct {
	Length  int `json:"length"`
	Type    int `json:"type"`
	Content any `json:"content"`
}

// 内容包flag类型
const (
	HANDS = iota
	HEARTBEAT
	DATA
	DISCONNECT
)

// 内容包(req)
type ContentPackReq struct {
	Flag  int    `json:"flag"`
	Route string `json:"route"`
	Data  any    `json:"data"`
}

// 内容包(res)
type ContentPackRes struct {
	Flag    int    `json:"flag"`
	Route   string `json:"route"`
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data"`
}

// 任务包(req)
type TaskPackReq struct {
	MessageType      int
	PacketFromConnId int64
	ContentFlag      int
	ContentRoute     string
	ContentData      any
}

// 任务包(res)
type TaskPackRes struct {
	MessageType       int
	PacketSendConnIds []int64
	ContentFlag       int
	ContentRoute      string
	ContentCode       int
	ContentMessage    string
	ContentData       any
}

// ws消息包解码
func WSDecodePacket(bytes []byte) (*MessagePacket, error) {
	messagePacket := &MessagePacket{}
	err := json.Unmarshal(bytes, messagePacket)
	if err != nil {
		return nil, err
	}
	return messagePacket, nil
}

// ws消息包编码
func WSEncodePacket(message *MessagePacket) ([]byte, error) {
	contentBytes, err := json.Marshal(message.Content)
	if err != nil {
		return nil, err
	}
	message.Length = len(contentBytes)
	messageBytes, err := json.Marshal(message)
	if err != nil {
		return nil, err
	}
	return messageBytes, nil
}

// tcp消息包解码
func TCPDecodePacket(conn *net.TCPConn) (*MessagePacket, error) {
	messagePacket := &MessagePacket{}

	lengthArr := make([]byte, 8)
	if _, err := conn.Read(lengthArr); err != nil {
		return nil, err
	}
	lengthBuff := bytes.NewReader(lengthArr)
	if err := binary.Read(lengthBuff, binary.LittleEndian, messagePacket.Length); err != nil {
		return nil, err
	}

	typeArr := make([]byte, 8)
	if _, err := conn.Read(typeArr); err != nil {
		return nil, err
	}
	typeBuff := bytes.NewReader(typeArr)
	if err := binary.Read(typeBuff, binary.LittleEndian, messagePacket.Type); err != nil {
		return nil, err
	}

	bodyArr := make([]byte, messagePacket.Length)
	if _, err := conn.Read(bodyArr); err != nil {
		return nil, err
	}
	bodyBuff := bytes.NewReader(bodyArr)
	if err := binary.Read(bodyBuff, binary.LittleEndian, &messagePacket.Content); err != nil {
		return nil, err
	}
	return messagePacket, nil
}

// tcp消息包编码
func TCPEncodePacket(message *MessagePacket) ([]byte, error) {
	contentBytes, err := convertor.ToBytes(message.Content)
	if err != nil {
		return nil, err
	}
	contentLength := len(contentBytes)
	message.Length = contentLength

	dataBuff := bytes.NewBuffer([]byte{})
	if err := binary.Write(dataBuff, binary.LittleEndian, message.Length); err != nil {
		return nil, err
	}
	if err := binary.Write(dataBuff, binary.LittleEndian, message.Type); err != nil {
		return nil, err
	}
	if err := binary.Write(dataBuff, binary.LittleEndian, message.Content); err != nil {
		return nil, err
	}
	return dataBuff.Bytes(), nil
}
