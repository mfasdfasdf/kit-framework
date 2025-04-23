package health

import (
	"framework/log"
	"github.com/arl/statsviz"
	"net/http"
	"strconv"
)

func RunServer(port int) {
	mux := http.NewServeMux()
	if err := statsviz.Register(mux); err != nil {
		log.Fatal("======>监控注册失败!======>err:", err)
	}
	if err := http.ListenAndServe(":"+strconv.Itoa(port), mux); err != nil {
		log.Fatal("======>监控启动失败!======>err:", err)
	}
}
