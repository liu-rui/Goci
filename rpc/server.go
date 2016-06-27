package rpc

import (
	"errors"

	"git.apache.org/thrift.git/lib/go/thrift"
)

var (
	errorServerNoInit = errors.New("please execute Serve function  first")
)

//Server PRC服务端类型
type Server struct {
	address      string
	processor    thrift.TProcessor
	simpleServer *thrift.TSimpleServer
}

//Serve 启动
func (server *Server) Serve() error {
	transportFactory := thrift.NewTTransportFactory()
	protocolFactory := thrift.NewTCompactProtocolFactory()
	serverTransport, err := thrift.NewTServerSocket(server.address)

	if err != nil {
		return err
	}
	server.simpleServer = thrift.NewTSimpleServer4(server.processor, serverTransport, transportFactory, protocolFactory)

	if err := server.simpleServer.Serve(); err != nil {
		return err
	}
	return nil
}

//Stop 停止
func (server *Server) Stop() error {
	if server.simpleServer == nil {
		return errorServerNoInit
	}

	return server.simpleServer.Stop()
}

//NewServer 创建PRC服务端对象
func NewServer(address string, processor thrift.TProcessor) *Server {
	return &Server{address, processor, nil}
}
