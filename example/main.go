package main

import (
	"example/application"
	"time"

	rclient "config-platform/client"
	"config-platform/resolver"

	clientv3 "go.etcd.io/etcd/client/v3"
)

type testClientConn struct {
}

func (c testClientConn) UpdateState(addr string) error {
	return nil
}

func main() {
	etcdCli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"localhost:2379"},
		DialTimeout: 3 * time.Second,
	})
	if err != nil {
		panic(err)
	}
	cc := testClientConn{}
	builder := resolver.NewBuilder(etcdCli)
	etcdResolver, err := builder.Build("test", cc)
	if err != nil {
		panic(err)
	}

	discoverCli := rclient.NewClient(rclient.ClientConfig{}, rclient.WithDiscovery(etcdResolver))
	app := application.NewAppForTest(discoverCli)
	app.Start()
}
