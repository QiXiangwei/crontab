package main

import (
	"crontab/master"
	"flag"
	"runtime"
)

var (
	confFile string
)

func initArgs() {
	flag.StringVar(&confFile, "config", "config.json", "")
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

	if err = master.InitConfig(confFile); err != nil {
		return
	}

	if err = master.InitApiServer(); err != nil {
		return
	}
}
