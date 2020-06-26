package config

import (
	"encoding/json"
	"io/ioutil"
)

type MasterConfig struct {
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
	GMasterConfig *MasterConfig
)

func InitMasterConfig(confFile string) (err error) {
	var (
		context []byte
		config  *MasterConfig
	)

	if context, err = ioutil.ReadFile(confFile); err != nil {
		return
	}

	if err = json.Unmarshal(context, &config); err != nil {
		return
	}

	GMasterConfig = config
	return
}
