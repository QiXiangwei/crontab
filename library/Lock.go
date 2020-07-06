package library

import (
	"context"
	"crontab/common"
	"errors"
	"github.com/coreos/etcd/clientv3"
)

type Lock struct {
	kv         clientv3.KV
	lease      clientv3.Lease
	jobName    string
	cancelFunc context.CancelFunc
	leaseId    clientv3.LeaseID
	isLocked   bool
}

func InitLock(jobName string, kv clientv3.KV, lease clientv3.Lease) (lock *Lock) {
	lock = &Lock{
		kv:           kv,
		lease:        lease,
		jobName:      jobName,
	}
	return
}

func (lock *Lock) TryLock() (err error) {
	var (
		leaseGrantResp *clientv3.LeaseGrantResponse
		leaseId        clientv3.LeaseID
		cancelCtx      context.Context
		cancelFunc     context.CancelFunc
		keepChanReps   <- chan *clientv3.LeaseKeepAliveResponse
		txn            clientv3.Txn
		txnResp        *clientv3.TxnResponse
		lockKey        string
	)

	if leaseGrantResp, err = lock.lease.Grant(context.TODO(), 5); err != nil {
		return
	}
	leaseId = leaseGrantResp.ID

	cancelCtx, cancelFunc = context.WithCancel(context.TODO())

	go func() {
		var (
			keepResp *clientv3.LeaseKeepAliveResponse
		)
		for {
			select {
			case keepResp = <- keepChanReps:
				if keepResp ==  nil {
					goto END
				}
			}
		}
		END:
	}()

	if keepChanReps, err = lock.lease.KeepAlive(context.TODO(), leaseId); err != nil {
		goto FAIL
	}

	txn     = lock.kv.Txn(context.TODO())
	lockKey = common.CRON_LOCK_KEY + lock.jobName
	txn.If(clientv3.Compare(clientv3.CreateRevision(lockKey), "=", 0)).
		Then(clientv3.OpPut(lockKey, "", clientv3.WithLease(leaseId))).
		Else(clientv3.OpGet(lockKey))

	if txnResp, err = txn.Commit(); err != nil {
		goto FAIL
	}

	if !txnResp.Succeeded {
		err = errors.New("锁已被占用")
	}

	lock.leaseId    = leaseId
	lock.cancelFunc = cancelFunc
	lock.isLocked   = true
	return
FAIL:
	cancelFunc()
	lock.lease.Revoke(context.TODO(), leaseId)
	return
}

func (lock *Lock) UnLock() {
	if lock.isLocked {
		lock.cancelFunc()
		lock.lease.Revoke(context.TODO(), lock.leaseId)
	}
}