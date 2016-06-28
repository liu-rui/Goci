package client

import (
	"errors"

	"git.apache.org/thrift.git/lib/go/thrift"
)

func NewTTransport(address string) (thrift.TTransport, error) {
	transportFactory := thrift.NewTTransportFactory()
	socket, err := thrift.NewTSocket(address)

	if err != nil {
		return nil, err
	}
	useTransport := transportFactory.GetTransport(socket)
	return useTransport, nil
}

func NewTProtocolFactory() thrift.TProtocolFactory {
	return thrift.NewTCompactProtocolFactory()
}

type ClientConfig struct {
	Direct        *Direct
	ServiceCentre *ServiceCentre
}

type Direct struct {
	Address string
}

type ServiceCentre struct {
	Servers []string
	Name    string
}

type Finder interface {
	Init() error
	Get() (string, error)
}

func NewFinder(conf *ClientConfig) (Finder, error) {
	if conf.Direct != nil {
		return newDirectFinder(conf.Direct), nil
	} else if conf.ServiceCentre != nil {
		return newServiceCentreFinder(conf.ServiceCentre), nil
	} else {
		return nil, errors.New("客户端配置信息无效，Direct或ServiceCentre必须配置一个")
	}
}
