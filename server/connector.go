package server

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"github.com/mfasdfasdf/kit-framework/log"
	"sync"
)

const (
	START = iota
	STOP
)

type IConnector interface {
	GetId() int64
	reader()
	writer()
	ReceivePacket(message *MessageResPacket)
}

//WS连接者

type WSConnector struct {
	RWLock     sync.RWMutex
	Id         int64
	Kind       string
	Conn       *websocket.Conn
	ReaderChan chan *MessageReqPacket
	WriterChan chan *MessageResPacket
	Status     int
}

func InitWSConnector(conn *websocket.Conn, id int64) *WSConnector {
	connector := &WSConnector{
		Id:         id,
		Conn:       conn,
		ReaderChan: make(chan *MessageReqPacket, 10),
		WriterChan: make(chan *MessageResPacket, 10),
		Status:     START,
	}
	go connector.reader()
	go connector.writer()
	return connector
}

func (o *WSConnector) GetId() int64 {
	return o.Id
}

func (o *WSConnector) reader() {
	for {
		_, bytes, err := o.Conn.ReadMessage()
		if err != nil {
			_manager.DelConn(o.Id)
			return
		}
		reqPacket, err := WsDecodePacket(bytes)
		if err != nil {
			log.Error("Decode packet err:%v", err)
			continue
		}
		//构建消息任务
		beforeTask := &TaskReq{
			PacketType:       reqPacket.Type,
			PacketFromConnId: o.Id,
			BodyFlat:         reqPacket.Body.Flag,
			BodyRoute:        reqPacket.Body.Route,
			BodyData:         reqPacket.Body.Data,
		}
		_dispatcher.receiveTask(beforeTask)
	}
}

func (o *WSConnector) writer() {
	for {
		select {
		case messageResPacket, ok := <-o.WriterChan:
			if !ok {
				return
			}
			bytes, err := json.Marshal(messageResPacket)
			if err != nil {
				return
			}
			log.Info("writerChan->messageResPacket:%v", string(bytes))
			resPacket, err := WsEncodePacket(messageResPacket)
			if err != nil {
				log.Error("EncodePacket err:%v", err)
			}
			err = o.Conn.WriteMessage(websocket.TextMessage, resPacket)
			if err != nil {
				log.Error("WriteMessage err:%v", err)
			}
		}
	}
}

func (o *WSConnector) ReceivePacket(message *MessageResPacket) { o.WriterChan <- message }

//TCP连接者

//type TCPConnector struct {
//	RWLock     sync.RWMutex
//	Id         int64
//	Kind       string
//	Conn       *net.TCPConn
//	ReaderChan chan *MessageReqPacket
//	WriterChan chan *MessageResPacket
//	Status     int
//}
//
//func InitTCPConnector(conn *net.TCPConn, id int64) *TCPConnector {
//	connector := &TCPConnector{
//		Id:         id,
//		Conn:       conn,
//		ReaderChan: make(chan *MessageReqPacket, 10),
//		WriterChan: make(chan *MessageResPacket, 10),
//		Status:     START,
//	}
//	go connector.reader()
//	go connector.writer()
//	return connector
//}
//
//func (o *TCPConnector) GetId() int64 {
//	return o.Id
//}
//
//func (o *TCPConnector) reader() {
//	for {
//		_, bytes, err := o.Conn.Read()
//		if err != nil {
//			_manager.DelConn(o.Id)
//			return
//		}
//		reqPacket, err := WsDecodePacket(bytes)
//		if err != nil {
//			log.Error("Decode packet err:%v", err)
//			continue
//		}
//		//解压的包添加请求者连接id
//		reqPacket.Body.FromConnId = o.Id
//		bytes, err = json.Marshal(reqPacket)
//		if err != nil {
//			continue
//		}
//		log.Info("readerChan<-messageReqPacket:%v", string(bytes))
//		//构建消息任务
//		beforeTask := &MessageReqTask{
//			PacketType:     reqPacket.Type,
//			BodyFlat:       reqPacket.Body.Flag,
//			BodyRoute:      reqPacket.Body.Route,
//			BodyData:       reqPacket.Body.Data,
//			BodyFromConnId: reqPacket.Body.FromConnId,
//		}
//		_dispatcher.receiveTask(beforeTask)
//	}
//}
//
//func (o *TCPConnector) writer() {
//	for {
//		select {
//		case messageResPacket, ok := <-o.WriterChan:
//			if !ok {
//				return
//			}
//			bytes, err := json.Marshal(messageResPacket)
//			if err != nil {
//				return
//			}
//			log.Info("writerChan->messageResPacket:%v", string(bytes))
//			resPacket, err := WsEncodePacket(messageResPacket)
//			if err != nil {
//				log.Error("EncodePacket err:%v", err)
//			}
//			err = o.Conn.WriteMessage(websocket.TextMessage, resPacket)
//			if err != nil {
//				log.Error("WriteMessage err:%v", err)
//			}
//		}
//	}
//}
//
//func (o *TCPConnector) ReceivePacket(message *MessageResPacket) { o.WriterChan <- message }
