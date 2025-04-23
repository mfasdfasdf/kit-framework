package etcd

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/duke-git/lancet/v2/convertor"
	"github.com/mfasdfasdf/kit-framework/config"
	"github.com/mfasdfasdf/kit-framework/log"
	"go.etcd.io/etcd/api/v3/mvccpb"
	"go.etcd.io/etcd/client/v3"
	"google.golang.org/grpc/resolver"
	"strconv"
	"time"
)

var EtcdCli *etcdClient

type etcdClient struct {
	cli         *clientv3.Client
	ServerInfos map[string][]string
	closeChan   chan bool
}

func InitEtcd() {
	if EtcdCli != nil {
		return
	}
	client, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{config.Configuration.Etcd.Url + ":" + strconv.Itoa(config.Configuration.Etcd.Port)},
		DialTimeout: time.Duration(config.Configuration.Etcd.DialTime) * time.Second,
	})
	if err != nil {
		log.Fatal("======>初始化etcd失败!<======")
	}
	EtcdCli = &etcdClient{cli: client, ServerInfos: make(map[string][]string), closeChan: make(chan bool, 1)}
	//向etcd注册信息
	EtcdCli.register()
}

func (e *etcdClient) register() {
	key := config.Configuration.AppName + "/" + config.Configuration.Version + "/" + convertor.ToString(config.Configuration.WorkId)
	value := make(map[string]any)
	value["grpcUrl"] = config.Configuration.Grpc.Url
	value["grpcPort"] = config.Configuration.Grpc.Port
	value["wsUrl"] = config.Configuration.Ws.Url
	value["wsPort"] = config.Configuration.Ws.Port
	marshal, err := json.Marshal(value)
	if err != nil {
		return
	}
	val := string(marshal)
	leaseId, err := e.SetLeaseId(int64(config.Configuration.Etcd.Ttl))
	if err != nil {
		return
	}
	e.DelKV(key)
	err = e.SetKV(key, val, leaseId)
	if err != nil {
		return
	}
	//开启心跳检测
	go e.heartbeat(leaseId)
}

func (e *etcdClient) SetLeaseId(ttl int64) (clientv3.LeaseID, error) {
	grant, err := EtcdCli.cli.Grant(context.Background(), ttl)
	if err != nil {
		log.Error("======>创建ETCD租约失败!======>err:", err.Error())
		return 0, err
	}
	return grant.ID, nil
}

func (e *etcdClient) DelLeaseId(leaseId clientv3.LeaseID) error {
	_, err := e.cli.Revoke(context.Background(), leaseId)
	if err != nil {
		return err
	}
	return nil
}

func (e *etcdClient) SetKV(key string, value string, leaseId clientv3.LeaseID) error {
	_, err := EtcdCli.cli.Put(context.Background(), key, value, clientv3.WithLease(leaseId))
	if err != nil {
		log.Warn("======>绑定ETCD租约失败!<======")
		return err
	}
	return nil
}

func (e *etcdClient) DelKV(key string) error {
	_, err := e.cli.Delete(context.Background(), key)
	if err != nil {
		return err
	}
	return nil
}

func (e *etcdClient) GetKV(key string) (string, error) {
	res, err := e.cli.Get(context.Background(), key)
	if err != nil {
		return "", err
	}
	if len(res.Kvs) == 0 {
		return "", errors.New(fmt.Sprintf("etcd key:%v 没有值!", key))
	}
	return string(res.Kvs[0].Value), nil
}

func (e *etcdClient) GetKVByPrefix(prefix string) ([]*mvccpb.KeyValue, error) {
	res, err := e.cli.Get(context.Background(), prefix, clientv3.WithPrefix())
	if err != nil {
		return nil, err
	}
	if len(res.Kvs) == 0 {
		return nil, errors.New(fmt.Sprintf("etcd prefix:%v 没有值!", prefix))
	}
	return res.Kvs, nil
}

func (e *etcdClient) RenewalLease(leaseId clientv3.LeaseID) {
	_, err := e.cli.KeepAliveOnce(context.Background(), leaseId)
	if err != nil {
		log.Error("======>续租ETCD失败! leaseId: %v<======", leaseId)
		return
	}
}

func (e *etcdClient) heartbeat(leaseId clientv3.LeaseID) {
	for {
		e.RenewalLease(leaseId)
		time.Sleep(5 * time.Second)
	}
}

func (e *etcdClient) Build(target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOptions) (resolver.Resolver, error) {
	log.Info("struct:[etcdClient]:::function:[Build]")
	prefix := target.URL.Host + target.URL.Path
	kvs, err := e.GetKVByPrefix(prefix)
	if err != nil {
		log.Error("======>grpc获取目标etcd失败!prefix: %v<======", prefix)
	}
	addresses := make([]resolver.Address, 0)
	for _, kv := range kvs {
		info := make(map[string]any)
		err = json.Unmarshal(kv.Value, &info)
		if err != nil {
			continue
		}
		addresses = append(addresses, resolver.Address{Addr: fmt.Sprintf("%v:%v", info["grpcUrl"], info["grpcPort"])})
	}
	err = cc.UpdateState(resolver.State{Addresses: addresses})
	if err != nil {
		log.Error("======>grpc更新地址失败! err: %v<======", err)
		return nil, err
	}
	return nil, nil
}

func (e *etcdClient) Scheme() string {
	log.Info("struct:[etcdClient]:::function:[Scheme]")
	return "etcd"
}
