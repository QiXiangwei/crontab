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
	etcClient  *clientv3.Client
	etcKv      clientv3.KV
	etcLease   clientv3.Lease
	etcWatcher clientv3.Watcher
}

var (
	GEtcServer *EtcServer
)

func InitMasterServer() (err error) {
	var(
		cli     *clientv3.Client
		kv      clientv3.KV
		lease   clientv3.Lease
		watcher clientv3.Watcher
	)

	if cli, err = clientv3.New(clientv3.Config{
		Endpoints:            config.GMasterConfig.EtcEndpoints,
		DialTimeout:          time.Duration(config.GMasterConfig.EtcDialTimeout) * time.Millisecond,
	}); err != nil {
		fmt.Println("etcd connect failed")
		return
	}

	kv      = clientv3.NewKV(cli)
	lease   = clientv3.NewLease(cli)
	watcher = clientv3.NewWatcher(cli)

	GEtcServer = &EtcServer{
		etcClient:  cli,
		etcKv:      kv,
		etcLease:   lease,
		etcWatcher: watcher,
	}

	return
}

func InitWorkerEtcServer() (err error) {
	var (
		cli     *clientv3.Client
		kv      clientv3.KV
		lease   clientv3.Lease
		watcher clientv3.Watcher
	)

	if cli, err = clientv3.New(clientv3.Config{
		Endpoints:            config.GWorkerConfig.EtcEndpoints,
		DialTimeout:          time.Duration(config.GWorkerConfig.EtcDialTimeout) * time.Millisecond,
	}); err != nil {
		return
	}

	kv      = clientv3.NewKV(cli)
	lease   = clientv3.NewLease(cli)
	watcher = clientv3.NewWatcher(cli)

	GEtcServer = &EtcServer{
		etcClient:  cli,
		etcKv:      kv,
		etcLease:   lease,
		etcWatcher: watcher,
	}

	if err = GEtcServer.watchJobs(common.CRON_JOB_KEY); err != nil {
		return
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

func (etcServer *EtcServer) CreateLease(leaseTime int64) (leaseId clientv3.LeaseID, err error) {
	var (
		leaseGrantResp *clientv3.LeaseGrantResponse
	)

	if leaseGrantResp, err = etcServer.etcLease.Grant(context.TODO(), leaseTime); err != nil {
		return
	}

	leaseId = leaseGrantResp.ID
	return
}

func (etcServer *EtcServer) Put(key string, value string) (result []byte, err error) {
	var (
		putResp *clientv3.PutResponse
	)

	if putResp, err = etcServer.etcKv.Put(context.TODO(), key, value, clientv3.WithPrevKV()); err != nil {
		return
	}

	result = putResp.PrevKv.Value
	return
}

func (etcServer *EtcServer) PutWithLease(key string, value string, leaseId clientv3.LeaseID) (err error) {

	if _, err = etcServer.etcKv.Put(context.TODO(), key, value, clientv3.WithLease(leaseId)); err != nil {
		return
	}
	return
}

func (etcServer *EtcServer) watchJobs(key string) (err error) {
	var (
		getResp            *clientv3.GetResponse
		kvPair             *mvccpb.KeyValue
		job                *common.Job
		watchStartRevision int64
		watchChan          clientv3.WatchChan
		watchResp          clientv3.WatchResponse
		watchEvent         *clientv3.Event
		jobEvent           *common.JobEvent
		jobName            string
	)

	if getResp, err = etcServer.etcKv.Get(context.TODO(), key, clientv3.WithPrefix()); err != nil {
		return
	}

	for _, kvPair = range getResp.Kvs {
		if job, err = common.UnmarshalJob(kvPair.Value); err == nil {
			jobEvent = common.BuildJobEvent(common.JOB_EVENT_SAVE, job)
			fmt.Println(*jobEvent)
		}
	}

	go func() {
		watchStartRevision = getResp.Header.Revision + 1

		watchChan = etcServer.etcWatcher.Watch(context.TODO(), key, clientv3.WithRev(watchStartRevision), clientv3.WithPrefix())

		for watchResp = range watchChan {
			for _, watchEvent = range watchResp.Events {
				switch watchEvent.Type {
				case mvccpb.PUT:
					if job, err = common.UnmarshalJob(watchEvent.Kv.Value); err != nil {
						continue
					}

					jobEvent = common.BuildJobEvent(common.JOB_EVENT_SAVE, job)
				case mvccpb.DELETE :
					jobName  = common.ExtractJobName(string(watchEvent.Kv.Key))
					job      = &common.Job{Name:jobName}
					jobEvent = common.BuildJobEvent(common.JOB_EVENT_DELETE, job)
				}
				fmt.Println(*jobEvent)
			}
		}
	}()

	return
}