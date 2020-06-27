package config

import (
	"encoding/json"
	"io/ioutil"
)

type WorkerConfig struct {
	EtcEndpoints    []string `json:"etc_endpoints"`
	EtcDialTimeout  int `json:"etc_dial_timeout"`
}

var (
	GWorkerConfig *WorkerConfig
)

func InitWorkerConfig(confFile string) (err error) {
	var (
		context []byte
		config  *WorkerConfig
	)

	if context, err = ioutil.ReadFile(confFile); err != nil {
		return
	}
	if err = json.Unmarshal(context, &config); err != nil {
		return
	}

	GWorkerConfig = config
	return
}
