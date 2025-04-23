package server

var Server *server

type server struct {
	Manager *manager
	Handler *handler
}

func InitServer() {
	if Server != nil {
		return
	}
	//初始化路由管理
	initHandler()
	//初始化分发者
	initDispatcher()
	//初始化管理者
	initManager()
	Server = &server{Manager: _manager, Handler: _handler}
}
