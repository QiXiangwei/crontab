package library

import (
	"github.com/coreos/etcd/clientv3"
)

type EtcServer struct {
	etcClient *clientv3.Client
	etcKv     clientv3.KV
	etcLease  clientv3.Lease
}

var (
	GEtcServer *EtcServer
)

func InitServer() (err error) {
	var(
		cli   *clientv3.Client
		kv    clientv3.KV
		lease clientv3.Lease
	)

	if cli, err = clientv3.New(clientv3.Config{
		Endpoints:            nil,
		DialTimeout:          0,
	}); err != nil {
		return
	}

	kv    = clientv3.NewKV(cli)
	lease = clientv3.NewLease(cli)

	GEtcServer = &EtcServer{
		etcClient: cli,
		etcKv:     kv,
		etcLease:  lease,
	}

	return
}
