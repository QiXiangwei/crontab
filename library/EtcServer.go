package library

import (
	"context"
	"crontab/common"
	"crontab/config"
	"encoding/json"
	"fmt"
	"github.com/coreos/etcd/clientv3"
	"time"
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
		Endpoints:            config.GConfig.EtcEndpoints,
		DialTimeout:          time.Duration(config.GConfig.EtcDialTimeout) * time.Millisecond,
	}); err != nil {
		fmt.Println("etcd connect failed")
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

func (etcServer *EtcServer) Save(job *common.Job) (oldJob *common.Job, err error) {
	var (
		jobValue []byte
		putRep   *clientv3.PutResponse
	)
	jobKey := "/cron/job/" + job.Name
	if jobValue, err = json.Marshal(job); err != nil {
		return
	}
	if putRep, err = etcServer.etcKv.Put(context.TODO(), jobKey, string(jobValue), clientv3.WithPrevKV()); err != nil {
		return
	}

	if putRep.PrevKv != nil {
		if err = json.Unmarshal(putRep.PrevKv.Value, &oldJob); err != nil {
			err = nil
			return
		}
	}

	return
}
