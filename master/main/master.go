package main

import (
	"crontab/config"
	"crontab/library"
	"crontab/master"
	"flag"
	"fmt"
	"runtime"
	"time"
)

var (
	confFile string
)

func initArgs() {
	flag.StringVar(&confFile, "../master_config.json", "worker_config.json", "")
	flag.Parse()
}

func initEnv() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}

func main()  {

	var (
		err error
	)

	initArgs()
	initEnv()

	if err = config.InitMasterConfig(confFile); err != nil {
		fmt.Println(err)
		return
	}


	if err = library.InitMasterServer(); err != nil {
		fmt.Println(err)
		return
	}


	if err = master.InitApiServer(); err != nil {
		fmt.Println(err)
		return
	}

	for {
		time.Sleep(1 * time.Second)
	}

}
