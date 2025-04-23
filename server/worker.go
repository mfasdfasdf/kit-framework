package server

import (
	"encoding/json"
	"github.com/mfasdfasdf/kit-framework/log"
)

type worker struct {
	workerId   int
	tasksQueue chan *MessageReqTask
}

func initWorker(workerId int, queueSize int) *worker {
	worker := &worker{
		workerId:   workerId,
		tasksQueue: make(chan *MessageReqTask, queueSize),
	}
	go worker.start()
	return worker
}

func (w *worker) QueueSize() int {
	return len(w.tasksQueue)
}

func (w *worker) pushTask(reqTask *MessageReqTask) {
	w.tasksQueue <- reqTask
}

func (w *worker) start() {
	for {
		select {
		case reqTask, ok := <-w.tasksQueue:
			if !ok {
				continue
			}
			bytes, err := json.Marshal(reqTask)
			if err != nil {
				continue
			}
			log.Info("======>workerId: %v, 处理任务:%v", w.workerId, string(bytes))
			handlerFunc := _handler.QueryHandler(reqTask.BodyRoute)
			if handlerFunc == nil {
				continue
			}
			resTask := handlerFunc(reqTask)

			if len(resTask.BodyToConnIds) == 0 {
				continue
			}
			connections := _manager.QueryBatchConn(resTask.BodyToConnIds)
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
