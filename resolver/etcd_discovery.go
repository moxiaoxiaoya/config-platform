package resolver

import (
	"context"
	"encoding/json"

	clientv3 "go.etcd.io/etcd/client/v3"
)

type builder struct {
	c *clientv3.Client
}

func NewBuilder(c *clientv3.Client) Builder {
	return builder{c: c}
}

func (b builder) Build(target string, cc ClientConn) (Discovery, error) {
	ctx, cancel := context.WithCancel(context.Background())
	d := &etcdDiscovery{
		target: target,
		cli:    b.c,
		ctx:    ctx,
		cancel: cancel,
	}

	upchan, err := d.watch()
	if err != nil {
		return nil, err
	}
	go d.updateAddr(upchan, cc)
	return d, nil
}

type etcdDiscovery struct {
	target string
	cli    *clientv3.Client
	ctx    context.Context
	cancel context.CancelFunc
}

func (e *etcdDiscovery) watch() (<-chan []*Update, error) {
	resp, err := e.cli.Get(e.ctx, e.target, clientv3.WithPrefix(), clientv3.WithSerializable())
	if err != nil {
		return nil, err
	}
	lg := e.cli.GetLogger()
	initUpdates := make([]*Update, 0, len(resp.Kvs))
	for _, kv := range resp.Kvs {
		var up Update
		if err := json.Unmarshal(kv.Value, &up); err != nil {
			lg.Warn("fail")
			continue
		}
		initUpdates = append(initUpdates, &up)
	}
	upch := make(chan []*Update, 1)
	if len(initUpdates) > 0 {
		upch <- initUpdates
	}
	go e.watchRealTime(upch)
	return upch, nil
}

func (e *etcdDiscovery) watchRealTime(upch chan []*Update) {
	defer close(upch)

	for {
		select {
		case <-e.ctx.Done():
			return
		}
	}
}

func (e *etcdDiscovery) updateAddr(upch <-chan []*Update, cc ClientConn) {
	for {
		select {
		case <-e.ctx.Done():
			return
		}
	}
}

func (e *etcdDiscovery) Discover() {

}

func (e *etcdDiscovery) Close() {
	e.cancel()
}
