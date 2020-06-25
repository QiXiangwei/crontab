package config

import (
	"encoding/json"
	"io/ioutil"
)

type Config struct {
	ApiPort         int `json:"api_port"`
	ApiReadTimeout  int `json:"api_read_timeout"`
	ApiWriteTimeOut int `json:"api_write_time_out"`

	EtcEndpoints    []string `json:"etc_endpoints"`
	EtcDialTimeout  int `json:"etc_dial_timeout"`

	RedisAddr       string `json:"redis_addr"`
	RedisPassword   string `json:"redis_password"`
	RedisDB         int `json:"redis_db"`

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