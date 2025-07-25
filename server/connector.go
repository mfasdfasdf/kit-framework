package server

import (
	"encoding/json"
	"github.com/duke-git/lancet/v2/convertor"
	"github.com/gorilla/websocket"
	"github.com/mfasdfasdf/kit-framework/log"
	"net"
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
	ReceiveMessagePacket(message *MessagePacket)
}

// WS连接者
type WSConnector struct {
	RWLock     sync.RWMutex
	Id         int64
	Kind       string
	Conn       *websocket.Conn
	ReaderChan chan *MessagePacket
	WriterChan chan *MessagePacket
	Status     int
}

func InitWSConnector(conn *websocket.Conn, id int64) *WSConnector {
	connector := &WSConnector{
		Id:         id,
		Conn:       conn,
		ReaderChan: make(chan *MessagePacket, 10),
		WriterChan: make(chan *MessagePacket, 10),
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
		_, messageBytes, err := o.Conn.ReadMessage()
		if err != nil {
			_manager.DelConn(o.Id)
			return
		}
		message, err := WSDecodePacket(messageBytes)
		if err != nil {
			log.Error("Decode packet err:%v", err)
			continue
		}
		messageJson, err := convertor.ToJson(message)
		if err != nil {
			log.Error("convert to json err:%v", err)
			continue
		}
		log.Info("=======>receive  message:%v", messageJson)
		//构建任务
		var contentPackReq ContentPackReq
		contentBytes, err := convertor.ToBytes(message.Content)
		if err != nil {
			log.Error("any to obj err:%v", err)
			continue
		}
		err = json.Unmarshal(contentBytes, &contentPackReq)
		if err != nil {
			log.Error("bytes to obj err:%v", err)
			continue
		}
		taskReq := &TaskPackReq{
			MessageType:      message.Type,
			PacketFromConnId: o.Id,
			ContentFlag:      contentPackReq.Flag,
			ContentRoute:     contentPackReq.Route,
			ContentData:      contentPackReq.Data,
		}
		_dispatcher.receiveTask(taskReq)
	}
}

func (o *WSConnector) writer() {
	for {
		select {
		case message, ok := <-o.WriterChan:
			if !ok {
				return
			}
			log.Info("writerChan -> message:%v", message)
			messageBytes, err := WSEncodePacket(message)
			if err != nil {
				log.Error("EncodePacket err:%v", err)
			}
			err = o.Conn.WriteMessage(websocket.TextMessage, messageBytes)
			if err != nil {
				log.Error("WriteMessage err:%v", err)
			}
		}
	}
}

func (o *WSConnector) ReceiveMessagePacket(message *MessagePacket) { o.WriterChan <- message }

// TCP连接者
type TCPConnector struct {
	RWLock     sync.RWMutex
	Id         int64
	Kind       string
	Conn       *net.TCPConn
	ReaderChan chan *MessagePacket
	WriterChan chan *MessagePacket
	Status     int
}

func InitTCPConnector(conn *net.TCPConn, id int64) *TCPConnector {
	connector := &TCPConnector{
		Id:         id,
		Conn:       conn,
		ReaderChan: make(chan *MessagePacket, 10),
		WriterChan: make(chan *MessagePacket, 10),
		Status:     START,
	}
	go connector.reader()
	go connector.writer()
	return connector
}

func (o *TCPConnector) GetId() int64 {
	return o.Id
}

func (o *TCPConnector) reader() {
	for {
		message, err := TCPDecodePacket(o.Conn)
		if err != nil {
			log.Error("DecodePacket err:%v", err)
			_manager.DelConn(o.Id)
			return
		}
		messageJson, _ := convertor.ToJson(message)
		log.Info("======>receive message:%v", messageJson)
		//构建任务
		var contentPackReq ContentPackReq
		contentBytes, err := convertor.ToBytes(message.Content)
		if err != nil {
			log.Error("any to obj err:%v", err)
			continue
		}
		err = json.Unmarshal(contentBytes, &contentPackReq)
		if err != nil {
			log.Error("bytes to obj err:%v", err)
			continue
		}
		taskReq := &TaskPackReq{
			MessageType:      message.Type,
			PacketFromConnId: o.Id,
			ContentFlag:      contentPackReq.Flag,
			ContentRoute:     contentPackReq.Route,
			ContentData:      contentPackReq.Data,
		}
		_dispatcher.receiveTask(taskReq)
	}
}

func (o *TCPConnector) writer() {
	for {
		select {
		case messageResPacket, ok := <-o.WriterChan:
			if !ok {
				return
			}
			bytesData, err := json.Marshal(messageResPacket)
			if err != nil {
				return
			}
			log.Info("======>send message:%v", string(bytesData))
			resBytes, err := TCPEncodePacket(messageResPacket)
			if err != nil {
				log.Error("EncodePacket err:%v", err)
			}
			_, err = o.Conn.Write(resBytes)
			if err != nil {
				log.Error("WriteMessage err:%v", err)
			}
		}
	}
}

func (o *TCPConnector) ReceiveMessagePacket(message *MessagePacket) { o.WriterChan <- message }
