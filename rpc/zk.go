package rpc

import (
	"time"

	"github.com/liu-rui/goci/log"
	"github.com/samuel/go-zookeeper/zk"
)

var (
	repairInterval = time.Second * 10
)

type ZKProcesser interface {
	Process() error
	ProcessWithEvent() error
}

type ZKObject struct {
	Servers   []string
	Conn      *zk.Conn
	Processer ZKProcesser
}

func (obj *ZKObject) Run() error {
	return obj.repairProcess()
}

func (obj *ZKObject) Close() {
	obj.closeConn()
}

func (obj *ZKObject) repairProcess() error {
	if obj.Conn != nil && obj.Conn.State() != zk.StateConnected {
		log.Debug("关闭连接")
		obj.closeConn()
	}

	if obj.Conn == nil {
		obj.createConn()
	}

	return obj.Processer.Process()
}

func (obj *ZKObject) Listen(watch <-chan zk.Event) {
	go func() {
		event := <-watch

		log.Debugf("节点发生更改 %v", event)
		switch event.State {
		case zk.StateDisconnected, zk.StateExpired, zk.StateUnknown, zk.StateAuthFailed:
			obj.startRepair()
			return
		}

		if err := obj.Processer.ProcessWithEvent(); err != nil {
			obj.startRepair()
		}
	}()
}

func (obj *ZKObject) startRepair() {
	go func() {
		log.Infof("与RPC服务中心%v断开连接，启动修复程序", obj.Servers)
		for {
			if err := obj.repairProcess(); err != nil {
				log.Error(err)
				time.Sleep(repairInterval)
			} else {
				break
			}
		}
		log.Infof("已与RPC服务中心%v重新建立连接，完成修复", obj.Servers)
	}()
}

func (obj *ZKObject) createConn() {
	c, _, err := zk.Connect(obj.Servers, time.Second*10)

	if err != nil {
		panic(err)
	}

	max := 5

	for c.State() != zk.StateConnected && max > 1 {
		time.Sleep(100 * time.Millisecond)
		max--
	}
	obj.Conn = c
}

func (obj *ZKObject) closeConn() {
	if obj.Conn == nil {
		return
	}

	obj.Conn.Close()
	obj.Conn = nil
}
