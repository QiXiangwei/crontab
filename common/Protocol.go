package common

import "encoding/json"

type Job struct {
	Name     string `json:"name"`
	Command  string `json:"command"`
	CronExpr string `json:"cronExpr"`
}

type Response struct {
	ErrNo  int `json:"errNo"`
	ErrStr string `json:"errStr"`
	Data   interface{} `json:"data"`
}

func BuildResponse(errNo int, errStr string, data interface{}) (resp []byte, err error) {
	var (
		rep Response
	)
	rep.ErrNo  = errNo
	rep.ErrStr = errStr
	rep.Data   = data
	resp, err = json.Marshal(rep)
	return
}