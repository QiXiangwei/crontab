package library

import (
	"context"
	"crontab/common"
	"crontab/config"
	"encoding/json"
	"fmt"
	"github.com/coreos/etcd/clientv3"
	"github.com/coreos/etcd/mvcc/mvccpb"
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
	jobKey := common.CRON_JOB_KEY + job.Name
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

func (etcServer *EtcServer) Delete(key string) (deleteJob *common.Job, err error) {
	var (
		deleteRep *clientv3.DeleteResponse
	)
	if deleteRep, err = etcServer.etcKv.Delete(context.TODO(), key, clientv3.WithPrevKV()); err != nil {
		return
	}
	if len(deleteRep.PrevKvs) != 0 {
		if err = json.Unmarshal(deleteRep.PrevKvs[0].Value, &deleteJob); err != nil {
			err = nil
			return
		}
	}
	return
}

func (etcServer *EtcServer) List(dirKey string) (jobList []*common.Job, err error) {
	var (
		getRep *clientv3.GetResponse
		kvPair *mvccpb.KeyValue
		temp   *common.Job
	)

	jobList = make([]*common.Job, 0)

	if getRep, err = etcServer.etcKv.Get(context.TODO(), dirKey, clientv3.WithPrefix()); err != nil {
		return
	}

	for _, kvPair = range getRep.Kvs {
		temp = &common.Job{}
		if err = json.Unmarshal(kvPair.Value, temp); err != nil {
			err = nil
			continue
		}
		jobList = append(jobList, temp)
	}
	return
}
