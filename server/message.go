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
	Length  int64 `json:"length"`
	Type    int64 `json:"type"`
	Content any   `json:"content"`
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
	MessageType      int64
	PacketFromConnId int64
	ContentFlag      int
	ContentRoute     string
	ContentData      any
}

// 任务包(res)
type TaskPackRes struct {
	MessageType       int64
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
	message.Length = int64(len(contentBytes))
	messageBytes, err := json.Marshal(message)
	if err != nil {
		return nil, err
	}
	return messageBytes, nil
}

// tcp消息包解码
func TCPDecodePacket(conn *net.TCPConn) (*MessagePacket, error) {
	message := &MessagePacket{}

	lengthArr := make([]byte, 8)
	if _, err := conn.Read(lengthArr); err != nil {
		return nil, err
	}
	if err := binary.Read(bytes.NewReader(lengthArr), binary.LittleEndian, &message.Length); err != nil {
		return nil, err
	}

	typeArr := make([]byte, 8)
	if _, err := conn.Read(typeArr); err != nil {
		return nil, err
	}
	if err := binary.Read(bytes.NewReader(typeArr), binary.LittleEndian, &message.Type); err != nil {
		return nil, err
	}

	contentArr := make([]byte, message.Length)
	if _, err := conn.Read(contentArr); err != nil {
		return nil, err
	}
	contentBytes := make([]byte, message.Length)
	if err := binary.Read(bytes.NewReader(contentArr), binary.LittleEndian, contentBytes); err != nil {
		return nil, err
	}
	err := json.Unmarshal(contentBytes, &message.Content)
	if err != nil {
		return nil, err
	}
	return message, nil
}

// tcp消息包编码
func TCPEncodePacket(message *MessagePacket) ([]byte, error) {
	contentBytes, err := convertor.ToBytes(message.Content)
	if err != nil {
		return nil, err
	}
	message.Length = int64(len(contentBytes))
	dataBuff := bytes.NewBuffer([]byte{})
	if err := binary.Write(dataBuff, binary.LittleEndian, message.Length); err != nil {
		return nil, err
	}
	if err := binary.Write(dataBuff, binary.LittleEndian, message.Type); err != nil {
		return nil, err
	}
	if err := binary.Write(dataBuff, binary.LittleEndian, contentBytes); err != nil {
		return nil, err
	}
	return dataBuff.Bytes(), nil
}
