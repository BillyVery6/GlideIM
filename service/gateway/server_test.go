package gateway

import (
	"go_im/pkg/rpc"
	"go_im/service"
	"testing"
)

func TestNewServer(t *testing.T) {

	config, err := service.GetConfig()
	if err != nil {
		panic(err)
	}
	etcd := config.Etcd.Servers

	server := NewServer(&rpc.ServerOptions{
		Name:        config.Gateway.Server.Name,
		Network:     config.Gateway.Server.Network,
		Addr:        config.Gateway.Server.Addr,
		Port:        config.Gateway.Server.Port,
		EtcdServers: etcd,
	})

	err = server.Run()
}
