package resolver

import (
	"context"
	"errors"

	clientv3 "go.etcd.io/etcd/client/v3"
)

const (
	_DefaultLease int64 = 500
)

var (
	ErrInvalidOperationType = errors.New("invalid operation type")
)

type OperationType uint8

const (
	Invalid OperationType = iota
	Add
	Delete
)

type ETCDRegister struct {
	lease int64
	cli   *clientv3.Client
}

func NewETCDRegister(cli *clientv3.Client, opts ...Option) *ETCDRegister {

	r := &ETCDRegister{cli: cli, lease: _DefaultLease}
	for _, opt := range opts {
		opt(r)
	}
	return r
}

type Option func(*ETCDRegister)

func WithLeashTime(lease int64) Option {
	return func(e *ETCDRegister) {
		e.lease = lease
	}
}

func (r *ETCDRegister) Register(ctx context.Context, key, addr string) error {
	resp, err := r.cli.Grant(ctx, r.lease)
	if err != nil {
		return err
	}
	up := UpdateOperation{
		Endpoint: Endpoint{Key: key, Addr: addr},
		opType:   Add,
		options:  []clientv3.OpOption{clientv3.WithLease(resp.ID)},
	}
	if err = r.Update(ctx, up); err != nil {
		return err
	}
	leaseChan, err := r.cli.KeepAlive(ctx, resp.ID)
	// TODO: how to deal with the leaseChan
	go func() { <-leaseChan }()
	return nil
}

func (r *ETCDRegister) LogOut(ctx context.Context, key, addr string) error {
	up := UpdateOperation{
		Endpoint: Endpoint{Key: key, Addr: addr},
		opType:   Delete}

	if err := r.Update(ctx, up); err != nil {
		return err
	}
	return nil
}

func (r *ETCDRegister) Update(ctx context.Context, up UpdateOperation) (err error) {
	var opReq clientv3.Op
	switch up.opType {
	case Add:
		// TODO: check key valid
		opReq = clientv3.OpPut(up.Key, up.Addr, up.options...)
	case Delete:
		opReq = clientv3.OpDelete(up.Key, up.options...)
	default:
		return ErrInvalidOperationType
	}
	_, err = r.cli.KV.Txn(ctx).Then(opReq).Commit()
	return err
}
