package master

import (
	"encoding/json"
	"io/ioutil"
)

type Config struct {
	apiPort         int `json:"api_port"`
	apiReadTimeout  int `json:"api_read_timeout"`
	apiWriteTimeOut int `json:"api_write_time_out"`
}

var (
	GConfig *Config
)

func InitConfig(confFile string) (err error) {
	var (
		context []byte
		config  *Config
	)

	if context, err = ioutil.ReadFile(confFile); err != nil {
		return
	}

	if err = json.Unmarshal(context, &config); err != nil {
		return
	}

	GConfig = config
	return
}