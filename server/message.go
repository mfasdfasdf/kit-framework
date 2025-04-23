package server

import "encoding/json"

// 消息包Type类型
const (
	HANDS = iota
	HEARTBEAT
	DATA
	DISCONNECT
)

type MessageReqPacket struct {
	Type   int            `json:"type"`
	Length int            `json:"length"`
	Body   MessageReqBody `json:"body"`
}

type MessageResPacket struct {
	Type   int            `json:"type"`
	Length int            `json:"length"`
	Body   MessageResBody `json:"body"`
}

const (
	ERROR   = 1000
	SUCCESS = 2000
	FAILURE = 5000
)

// 消息体Flag类型
const (
	REQUEST = iota
	RESPONSE
	NOTIFY
	PUSH
)

type MessageReqBody struct {
	Flag       int    `json:"flag"`
	Route      string `json:"route"`
	Data       any    `json:"data"`
	FromConnId int64  `json:"fromConnId"`
}

type MessageResBody struct {
	Flag    int    `json:"flag"`
	Route   string `json:"route"`
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data"`
}

// WsDecodePacket Ws消息包解码
func WsDecodePacket(bytes []byte) (*MessageReqPacket, error) {
	messagePacket := &MessageReqPacket{}
	err := json.Unmarshal(bytes, messagePacket)
	if err != nil {
		return nil, err
	}
	return messagePacket, nil
}

// WsEncodePacket Ws消息包编码
func WsEncodePacket(messageResPacket *MessageResPacket) ([]byte, error) {
	bodyBytes, err := json.Marshal(messageResPacket)
	if err != nil {
		return nil, err
	}
	messageResPacket.Length = len(bodyBytes)
	resBytes, err := json.Marshal(messageResPacket)
	if err != nil {
		return nil, err
	}
	return resBytes, nil
}

type MessageReqTask struct {
	PacketType     int
	BodyFlat       int
	BodyRoute      string
	BodyData       any
	BodyFromConnId int64
}

type MessageResTask struct {
	PacketType    int
	BodyFlat      int
	BodyRoute     string
	BodyCode      int
	BodyMessage   string
	BodyData      any
	BodyToConnIds []int64
}
