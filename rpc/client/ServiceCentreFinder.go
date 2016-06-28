package client

import (
	"fmt"
	"sync"

	"github.com/liu-rui/goci/log"
	"github.com/liu-rui/goci/rpc"
)

type serviceCentreFinder struct {
	lock     *sync.Mutex
	zkFinder *zkFinder
	servers  []string
	index    int
}

func (finder *serviceCentreFinder) Init() error {
	return finder.zkFinder.Run()
}

func (finder *serviceCentreFinder) Get() (string, error) {
	finder.lock.Lock()
	defer finder.lock.Unlock()

	if len(finder.servers) == 0 {
		return "", fmt.Errorf("目前在仓库%v上，没有针对%s的可用服务", finder.zkFinder.Servers, finder.zkFinder.path)
	}

	if finder.index == len(finder.servers) {
		finder.index = 0
	}
	server := finder.servers[finder.index]
	finder.index++
	return server, nil
}

func (finder *serviceCentreFinder) receive(servers []string) {
	if len(servers) == 0 {
		return
	}
	finder.lock.Lock()
	defer finder.lock.Unlock()

	finder.servers = servers
}

func newServiceCentreFinder(conf *ServiceCentre) *serviceCentreFinder {
	ret := &serviceCentreFinder{lock: new(sync.Mutex)}
	ret.zkFinder = newZKFinder(conf.Servers, conf.Name, ret)

	return ret
}

type serverChangedReceiver interface {
	receive([]string)
}

type zkFinder struct {
	rpc.ZKObject
	path                  string
	serverChangedReceiver serverChangedReceiver
}

func (finder *zkFinder) Process() error {
	return finder.ProcessWithEvent()
}

func (finder *zkFinder) ProcessWithEvent() error {
	servers, _, watch, err := finder.Conn.ChildrenW(finder.path)

	if err != nil {
		return err
	}
	log.Infof("RPC路径%s发现有新的服务器列表,服务器列表为：%v", finder.path, servers)

	finder.serverChangedReceiver.receive(servers)
	finder.Listen(watch)
	return nil
}

func newZKFinder(servers []string, path string, serverChangedReceiver serverChangedReceiver) *zkFinder {
	ret := &zkFinder{ZKObject: rpc.ZKObject{Servers: servers}, path: path, serverChangedReceiver: serverChangedReceiver}
	ret.Processer = ret
	return ret
}
