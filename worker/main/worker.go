package main

import (
	"crontab/config"
	"crontab/library"
	"crontab/worker"
	"flag"
	"runtime"
	"time"
)

var (
	confFile string
)

func initArgs() {
	flag.StringVar(&confFile, "worker_config", "./worker/worker_config.json", "")
	flag.Parse()
}

func initEnv() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}

func main() {
	var (
		err error
	)

	initArgs()
	initEnv()

	if err = config.InitWorkerConfig(confFile); err != nil {
		return
	}

	if err = worker.InitExecutor(); err != nil {
		return
	}

	if err = worker.InitScheduler(); err != nil {
		return
	}

	if err = library.InitWorkerEtcServer(); err != nil {
		return
	}

	for {
		time.Sleep(1 * time.Second)
	}
}