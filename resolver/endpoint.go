package resolver

import (
	clientv3 "go.etcd.io/etcd/client/v3"
)

type Endpoint struct {
	Key  string
	Addr string
}

type UpdateOperation struct {
	Endpoint
	opType  OperationType
	options []clientv3.OpOption
}

type Update struct {
	opType OperationType
	Key    string
	Addr   string
}
