package master

import (
	"crontab/common"
	"crontab/config"
	"crontab/library"
	"encoding/json"
	"github.com/coreos/etcd/clientv3"
	"net"
	"net/http"
	"strconv"
	"time"
)

type ApiServer struct {
	HttpServer *http.Server
}

var (
	GApiServer *ApiServer
)

func InitApiServer() (err error) {
	var (
		httpServer *http.Server
		mux        *http.ServeMux
		listen     net.Listener
	)

	mux = http.NewServeMux()
	mux.HandleFunc("/job/save", handleJobSave)
	mux.HandleFunc("/job/delete", handleJobDelete)
	mux.HandleFunc("/job/list", handleJobList)
	mux.HandleFunc("/job/kill", handleJobKill)

	mux.HandleFunc("/job/result", handleJobResult)

	if listen, err = net.Listen("tcp", ":" + strconv.Itoa(config.GMasterConfig.ApiPort)); err != nil {
		return
	}

	httpServer = &http.Server{
		Handler:           mux,
		ReadTimeout:       time.Duration(config.GMasterConfig.ApiReadTimeout) * time.Millisecond,
		WriteTimeout:      time.Duration(config.GMasterConfig.ApiWriteTimeOut) * time.Millisecond,
	}

	GApiServer = &ApiServer{HttpServer:httpServer}

	go httpServer.Serve(listen)

	return
}

func handleJobResult(rep http.ResponseWriter, req *http.Request) {
	var (
		err       error
		data      string
		jobName   string
		key       string
		resut     []byte
		jobResult *common.JobExecuteResult

	)

	if err = req.ParseForm(); err != nil {
		goto ERR
	}

	jobName = req.PostForm.Get("jobName")
	key     = common.REDIS_CRON_RESULT + jobName
	if data, err = library.GRedisServer.GetCacheData(key); err != nil {
		goto ERR
	}

	if err = json.Unmarshal([]byte(data), jobResult); err != nil {
		goto ERR
	}

	if resut, err = common.BuildResponse(0, "success", jobResult); err == nil {
		rep.Write(resut)
	}

	return
ERR:
	if resut, err = common.BuildResponse(-1, err.Error(), nil); err == nil {
		rep.Write(resut)
	}
}

func handleJobKill(rep http.ResponseWriter, req *http.Request) {
	var (
		err     error
		result  []byte
		name    string
		key     string
		leaseId clientv3.LeaseID
	)


	if err = req.ParseForm(); err != nil {
		goto ERR
	}

	name = req.PostForm.Get("killName")
	key  = common.CRON_KILL_KEY + name
	if leaseId, err = library.GEtcServer.CreateLease(1); err != nil {
		goto ERR
	}

	if err = library.GEtcServer.PutWithLease(key, "", leaseId); err != nil {
		goto ERR
	}

	if result, err = common.BuildResponse(0, "success", nil); err == nil {
		rep.Write(result)
	}

	return
ERR:
	if result, err = common.BuildResponse(-1, err.Error(), nil); err == nil {
		rep.Write(result)
	}
}

func handleJobList(rep http.ResponseWriter, req *http.Request) {
	var (
		err error
		dirKey string
		jobList []*common.Job
		result  []byte
	)

	dirKey = common.CRON_JOB_KEY
	if jobList, err = library.GEtcServer.List(dirKey); err != nil {
		goto ERR
	}

	if result, err = common.BuildResponse(0, "success", jobList); err == nil {
		rep.Write(result)
	}

	return
ERR:
	if result, err = common.BuildResponse(-1, err.Error(), nil); err == nil {
		rep.Write(result)
	}
}


func handleJobSave(rep http.ResponseWriter, req *http.Request) {
	var (
		err    error
		job    *common.Job
		old    *common.Job
		result []byte
		posJob string
	)
	if err = req.ParseForm(); err != nil {
		goto ERR
	}
	posJob = req.PostFormValue("job")
	if err = json.Unmarshal([]byte(posJob), &job); err != nil {
		goto ERR
	}

	if old, err = library.GEtcServer.Save(job); err != nil {
		goto ERR
	}
	if result, err = common.BuildResponse(0, "success", old); err == nil {
		rep.Write(result)
	}
	return
ERR:
	if result, err = common.BuildResponse(-1, err.Error(), nil); err == nil {
		rep.Write(result)
	}
}

func handleJobDelete(rep http.ResponseWriter, req *http.Request) {
	var (
		err        error
		result     []byte
		deleteName string
		deleteKey  string
		deleteJob  *common.Job
	)
	if err = req.ParseForm(); err != nil {
		goto ERR
	}

	deleteName = req.PostForm.Get("deleteName")
	deleteKey  = common.CRON_JOB_KEY + deleteName
	if deleteJob, err = library.GEtcServer.Delete(deleteKey); err != nil {
		goto ERR
	}

	if result, err = common.BuildResponse(0, "success", deleteJob); err == nil {
		rep.Write(result)
	}
	return

ERR:
	if result, err = common.BuildResponse(-1, err.Error(), nil); err == nil {
		rep.Write(result)
	}
}
