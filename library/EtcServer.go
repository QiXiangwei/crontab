package library

import (
	"github.com/coreos/etcd/clientv3"
)

type EtcServer struct {
	clientv3.Client
}
