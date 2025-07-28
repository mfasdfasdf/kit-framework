package server

import (
	"github.com/duke-git/lancet/v2/convertor"
	"github.com/mfasdfasdf/kit-framework/log"
)

type worker struct {
	workerId   int
	tasksQueue chan *TaskPackReq
}

func initWorker(workerId int, queueSize int) *worker {
	worker := &worker{
		workerId:   workerId,
		tasksQueue: make(chan *TaskPackReq, queueSize),
	}
	go worker.start()
	return worker
}

func (w *worker) QueueSize() int {
	return len(w.tasksQueue)
}

func (w *worker) pushTask(taskReq *TaskPackReq) {
	w.tasksQueue <- taskReq
}

func (w *worker) start() {
	for {
		select {
		case taskReq, ok := <-w.tasksQueue:
			if !ok {
				continue
			}
			taskReqJson, err := convertor.ToJson(taskReq)
			if err != nil {
				log.Error("DecodePacket err:%v", err)
				continue
			}
			log.Info("======>workerId: %v, 处理任务:%v", w.workerId, taskReqJson)
			handlerFunc := _handler.QueryHandler(taskReq.ContentRoute)
			if handlerFunc == nil {
				continue
			}
			taskRes := handlerFunc(taskReq)
			if taskRes == nil {
				continue
			}

			if len(taskRes.PacketSendConnIds) == 0 {
				continue
			}
			connections := _manager.QueryBatchConn(taskRes.PacketSendConnIds)
			for _, conn := range connections {
				messageResPacket := &MessagePacket{
					Type: taskRes.MessageType,
					Content: ContentPackRes{
						Flag:    taskRes.ContentFlag,
						Route:   taskRes.ContentRoute,
						Code:    taskRes.ContentCode,
						Message: taskRes.ContentMessage,
						Data:    taskRes.ContentData,
					},
				}
				(*conn).ReceiveMessagePacket(messageResPacket)
			}
		}
	}
}
