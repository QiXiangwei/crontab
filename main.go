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

	if err = config.InitConfig(confFile); err != nil {
		fmt.Println(err)
		return
	}


	if err = library.InitServer(); err != nil {
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