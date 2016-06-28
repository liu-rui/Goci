 package server
 
 

import (
	"errors"

	"git.apache.org/thrift.git/lib/go/thrift"
)

var (
	errorServerNoInit = errors.New("please execute Serve function  first")
)

//ServerConfig RPC服务端配置
type ServerConfig struct {
	Address       string
	ServiceCentre *ServiceCentre
}

//ServiceCentre 服务中心配置
type ServiceCentre struct {
	Servers []string
	Name    string
}

//Server PRC服务端类型
type Server struct {
	config       *ServerConfig
	simpleServer *thrift.TSimpleServer
}

//Serve 启动
func (server *Server) Serve(processor thrift.TProcessor) error {
	transportFactory := thrift.NewTTransportFactory()
	protocolFactory := thrift.NewTCompactProtocolFactory()
	serverTransport, err := thrift.NewTServerSocket(server.config.Address)

	if err != nil {
		return err
	}
	server.simpleServer = thrift.NewTSimpleServer4(processor, serverTransport, transportFactory, protocolFactory)

	if server.config.ServiceCentre != nil {
		serverPublisher := newServerPublisher(server.config.ServiceCentre.Servers, server.config.ServiceCentre.Name, server.config.Address)
		if err := serverPublisher.Run(); err != nil {
			panic(err)
		}
		defer serverPublisher.Close()
	}

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
func NewServer(address string) *Server {
	return &Server{config: &ServerConfig{Address: address}}
}

//NewServerByConfig  通过配置创建PRC服务端对象
func NewServerByConfig(conf *ServerConfig) *Server {
	return &Server{config: conf}
}
