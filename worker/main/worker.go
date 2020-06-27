package main

import (
	"crontab/config"
	"crontab/library"
	"flag"
	"fmt"
	"runtime"
	"time"
)

var (
	confFile string
)

func initArgs() {
	flag.StringVar(&confFile, "worker_config", "worker_config.json", "")
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
		fmt.Println(err.Error())
		return
	}

	if err = library.InitWorkerEtcServer(); err != nil {
		return
	}

	for {
		time.Sleep(1 * time.Second)
	}
}