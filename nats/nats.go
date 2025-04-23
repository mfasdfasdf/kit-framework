package nats

import (
	"fmt"
	"framework/config"
	"framework/log"
	"github.com/nats-io/nats.go"
)

var NatsClient *natsClient

type natsClient struct {
	cli *nats.Conn
}

func InitNatsClient() {
	if NatsClient != nil {
		return
	}
	conn, err := nats.Connect(fmt.Sprintf("%v:%v", config.Configuration.Nats.Url, config.Configuration.Nats.Port))
	if err != nil {
		return
	}
	NatsClient = &natsClient{cli: conn}
}

func (n *natsClient) Sub(name string) {
	_, err := NatsClient.cli.Subscribe(name, func(msg *nats.Msg) {
		log.Info("msg:%v", string(msg.Data))
	})
	if err != nil {
		log.Error("err:%v", err)
		return
	}
}

func (n *natsClient) Pub(name string, data []byte) {
	err := NatsClient.cli.Publish(name, data)
	if err != nil {
		log.Error("err:%v", err)
	}
}
