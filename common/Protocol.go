package common

import (
	"encoding/json"
	"github.com/gorhill/cronexpr"
	"strings"
	"time"
)

type Job struct {
	Name     string `json:"name"`
	Command  string `json:"command"`
	CronExpr string `json:"cronExpr"`
}

type JobSchedulerPlan struct {
	Job      *Job
	Expr     *cronexpr.Expression
	NextTime time.Time
}

type Response struct {
	ErrNo  int `json:"errNo"`
	ErrStr string `json:"errStr"`
	Data   interface{} `json:"data"`
}

type JobEvent struct {
	EventType int `json:"eventType"`
	Job       *Job
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

func UnmarshalJob(jobByte []byte) (job *Job, err error) {
	var (
		tempJob *Job
	)
	tempJob = &Job{}
	if err = json.Unmarshal(jobByte, &tempJob); err !=  nil  {
		return
	}
	job = tempJob
	return
}

func ExtractJobName(jobName string) (string) {
	return strings.TrimPrefix(jobName, CRON_JOB_KEY)
}

func BuildJobEvent(eventType int, job *Job) (jobEvent *JobEvent) {
	return &JobEvent{
		EventType: eventType,
		Job:       job,
	}
}

func BuildJobSchedulePlan(job *Job) (jobSchedulerPlan *JobSchedulerPlan, err error) {
	var (
		expr *cronexpr.Expression
	)

	if expr, err = cronexpr.Parse(job.CronExpr); err != nil {
		return
	}

	jobSchedulerPlan = &JobSchedulerPlan{
		Job:      job,
		Expr:     expr,
		NextTime: expr.Next(time.Now()),
	}

	return
}