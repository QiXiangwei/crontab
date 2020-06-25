package master

import (
	"crontab/common"
	"crontab/config"
	"crontab/library"
	"encoding/json"
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

	if listen, err = net.Listen("tcp", ":" + strconv.Itoa(config.GConfig.ApiPort)); err != nil {
		return
	}

	httpServer = &http.Server{
		Handler:           mux,
		ReadTimeout:       time.Duration(config.GConfig.ApiReadTimeout) * time.Millisecond,
		WriteTimeout:      time.Duration(config.GConfig.ApiWriteTimeOut) * time.Millisecond,
	}

	GApiServer = &ApiServer{HttpServer:httpServer}

	go httpServer.Serve(listen)

	return
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
