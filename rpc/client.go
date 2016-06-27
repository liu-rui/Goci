package rpc

import "git.apache.org/thrift.git/lib/go/thrift"

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
