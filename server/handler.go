package server

import "github.com/mfasdfasdf/kit-framework/log"

type LogicHandler func(task *TaskPackReq) *TaskPackRes

var _handler *handler

type handler struct {
	handlers map[string]LogicHandler
}

func initHandler() {
	_handler = &handler{handlers: make(map[string]LogicHandler)}
	log.Info("======>初始化路由管理完成!<======")
}

func (h *handler) AddHandler(key string, handler LogicHandler) {
	h.handlers[key] = handler
}

func (h *handler) QueryHandler(key string) LogicHandler {
	return h.handlers[key]
}
