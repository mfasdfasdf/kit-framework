package server

import (
	"github.com/mfasdfasdf/kit-framework/config"
	"github.com/mfasdfasdf/kit-framework/log"
	"sync"
)

var _dispatcher *dispatcher = nil

type dispatcher struct {
	lock            sync.RWMutex
	currentWorkerId int
	workers         []*worker
}

func initDispatcher() {
	if _dispatcher != nil {
		return
	}
	workerSize := config.Configuration.Distribute.WorkerSize
	if workerSize <= 0 {
		workerSize = 1
	}
	queueSize := config.Configuration.Distribute.QueueSize
	if queueSize <= 0 {
		queueSize = 1
	}
	workers := make([]*worker, 0)
	for i := 0; i < workerSize; i++ {
		worker := initWorker(i, queueSize)
		workers = append(workers, worker)
	}
	_dispatcher = &dispatcher{workers: workers, currentWorkerId: 0}
	log.Info("======>初始化分发管理完成!<======")
}

func (d *dispatcher) receiveTask(task *TaskPackReq) {
	d.lock.Lock()
	defer d.lock.Unlock()
	workerSize := len(d.workers)
	d.workers[d.currentWorkerId].pushTask(task)
	d.currentWorkerId++
	if d.currentWorkerId >= len(d.workers) {
		d.currentWorkerId = d.currentWorkerId % workerSize
	}
}
