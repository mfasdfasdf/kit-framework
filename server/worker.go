package server

import (
	"encoding/json"
	"github.com/mfasdfasdf/kit-framework/log"
)

type worker struct {
	workerId   int
	tasksQueue chan *TaskReq
}

func initWorker(workerId int, queueSize int) *worker {
	worker := &worker{
		workerId:   workerId,
		tasksQueue: make(chan *TaskReq, queueSize),
	}
	go worker.start()
	return worker
}

func (w *worker) QueueSize() int {
	return len(w.tasksQueue)
}

func (w *worker) pushTask(taskReq *TaskReq) {
	w.tasksQueue <- taskReq
}

func (w *worker) start() {
	for {
		select {
		case taskReq, ok := <-w.tasksQueue:
			if !ok {
				continue
			}
			bytes, err := json.Marshal(taskReq)
			if err != nil {
				continue
			}
			log.Info("======>workerId: %v, 处理任务:%v", w.workerId, string(bytes))
			handlerFunc := _handler.QueryHandler(taskReq.BodyRoute)
			if handlerFunc == nil {
				continue
			}
			resTask := handlerFunc(taskReq)

			if len(resTask.PacketToConnIds) == 0 {
				continue
			}
			connections := _manager.QueryBatchConn(resTask.PacketToConnIds)
			for _, conn := range connections {
				messageResPacket := &MessageResPacket{
					Type: resTask.PacketType,
					Body: MessageResBody{
						Flag:    resTask.BodyFlat,
						Route:   resTask.BodyRoute,
						Code:    resTask.BodyCode,
						Message: resTask.BodyMessage,
						Data:    resTask.BodyData,
					},
				}
				(*conn).ReceivePacket(messageResPacket)
			}
		}
	}
}
