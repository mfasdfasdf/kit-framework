package server

import "encoding/json"

const (
	ERROR   = 1000
	SUCCESS = 2000
	FAILURE = 5000
)

// 消息体type类型
const (
	REQUEST = iota
	RESPONSE
	NOTIFY
)

// 消息包flag类型
const (
	HANDS = iota
	HEARTBEAT
	DATA
	DISCONNECT
)

type MessageReqPacket struct {
	Length int            `json:"length"`
	Type   int            `json:"type"`
	Body   MessageReqBody `json:"body"`
}

type MessageReqBody struct {
	Flag  int    `json:"flag"`
	Route string `json:"route"`
	Data  any    `json:"data"`
}

type MessageResPacket struct {
	Length int            `json:"length"`
	Type   int            `json:"type"`
	Body   MessageResBody `json:"body"`
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

type TaskReq struct {
	PacketType       int
	PacketFromConnId int64
	BodyFlat         int
	BodyRoute        string
	BodyData         any
}

type TaskRes struct {
	PacketType      int
	PacketToConnIds []int64
	BodyFlat        int
	BodyRoute       string
	BodyCode        int
	BodyMessage     string
	BodyData        any
}
