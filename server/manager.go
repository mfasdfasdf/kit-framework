package server

import (
	"github.com/duke-git/lancet/v2/maputil"
	"github.com/mfasdfasdf/kit-framework/log"
	"sync"
)

const (
	CLOSE = iota
	OPEN
)

var _manager *manager = nil

type manager struct {
	lock        sync.RWMutex
	connections *maputil.ConcurrentMap[int64, *IConnector]
	status      int
}

func initManager() {
	if _manager != nil {
		return
	}
	_manager = &manager{
		connections: maputil.NewConcurrentMap[int64, *IConnector](0),
		status:      OPEN,
	}
	log.Info("======>初始化连接管理完成!<======")
}

func (m *manager) AddConn(conn IConnector) {
	m.connections.Set(conn.GetId(), &conn)
	size := m.ConnSize()
	log.Info("======>id:%v,已连接, 连接总数:%v<======", conn.GetId(), size)
}

func (m *manager) DelConn(id int64) {
	m.connections.Delete(id)
	size := m.ConnSize()
	log.Info("======>id:%v,断开, 连接总数:%v<======", id, size)
}

func (m *manager) ConnSize() int {
	total := 0
	m.connections.Range(func(k int64, v *IConnector) bool {
		total++
		return true
	})
	return total
}

func (m *manager) QueryOneConn(id int64) *IConnector {
	conn, has := m.connections.Get(id)
	if !has {
		return nil
	}
	return conn
}

func (m *manager) QueryBatchConn(ids []int64) []*IConnector {
	res := make([]*IConnector, 0)
	for _, id := range ids {
		conn, has := m.connections.Get(id)
		if !has {
			continue
		}
		if conn != nil {
			res = append(res, conn)
		}
	}
	return res
}
