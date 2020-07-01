package worker

import (
	"crontab/common"
	"crontab/library"
	"fmt"
	"time"
)

type Scheduler struct {
	JobEventChan      chan *common.JobEvent
	JobPlanTable      map[string]*common.JobSchedulerPlan
	JobExecutingTable map[string]*common.JobExecuteInfo
	JobExecuteChan    chan *common.JobExecuteResult
}

var (
	GScheduler *Scheduler
)

func (scheduler *Scheduler) TryStartJob(jobPlan *common.JobSchedulerPlan) {
	var (
		jobExecuteInfo *common.JobExecuteInfo
		jobExecuting    bool
	)

	if jobExecuteInfo, jobExecuting = scheduler.JobExecutingTable[jobPlan.Job.Name]; jobExecuting {
		fmt.Println("执行中:", jobExecuteInfo.Job.Name)
		return
	}

	jobExecuteInfo = common.BuildJobExecuteInfo(jobPlan)
	scheduler.JobExecutingTable[jobPlan.Job.Name] = jobExecuteInfo
	fmt.Println("准备执行:", jobExecuteInfo.Job.Name)
	GExecutor.ExecuteJob(jobExecuteInfo)
}

func (scheduler *Scheduler) TrySchedule() (scheduleAfter time.Duration) {
	var (
		jobPlan  *common.JobSchedulerPlan
		nowTime  time.Time
		nearTime *time.Time
	)

	if len(scheduler.JobPlanTable) == 0 {
		scheduleAfter = 1 * time.Second
		return
	}

	nowTime = time.Now()
	for _, jobPlan = range scheduler.JobPlanTable {
		if jobPlan.NextTime.Before(nowTime) || jobPlan.NextTime.Equal(nowTime) {
			scheduler.TryStartJob(jobPlan)
			jobPlan.NextTime = jobPlan.Expr.Next(nowTime)
		}
		if nearTime == nil || jobPlan.NextTime.Before(*nearTime) {
			nearTime = &jobPlan.NextTime
		}
	}
	scheduleAfter = (*nearTime).Sub(nowTime)
	return
}

func (scheduler *Scheduler) handleJobEvent(jobEvent *common.JobEvent) {
	var (
		jobSchedulePlan *common.JobSchedulerPlan
		err error
		jobExisted bool
	)
	switch jobEvent.EventType {
	case common.JOB_EVENT_SAVE:
		if jobSchedulePlan, err = common.BuildJobSchedulePlan(jobEvent.Job); err != nil {
			return
		}
		scheduler.JobPlanTable[jobEvent.Job.Name] = jobSchedulePlan
	case common.JOB_EVENT_DELETE:
		if jobSchedulePlan, jobExisted = scheduler.JobPlanTable[jobEvent.Job.Name]; jobExisted {
			delete(scheduler.JobPlanTable, jobEvent.Job.Name)
		}
		return
	}
}

func (scheduler *Scheduler) handleJobResult(jobResult *common.JobExecuteResult) {
	delete(scheduler.JobExecutingTable, jobResult.ExecuteInfo.Job.Name)
}

func (scheduler *Scheduler) schedulerLoop() {
	var (
		jobEvent      *common.JobEvent
		scheduleAfter time.Duration
		scheduleTimer *time.Timer
		jobResult     *common.JobExecuteResult
	)

	scheduleAfter = scheduler.TrySchedule()
	scheduleTimer = time.NewTimer(scheduleAfter)

	for {
		select {
		case jobEvent = <- scheduler.JobEventChan:
			scheduler.handleJobEvent(jobEvent)
		case <- scheduleTimer.C:
		case jobResult = <- scheduler.JobExecuteChan:
			fmt.Println("sign", jobResult.ExecuteInfo.Job.Name)
			scheduler.handleJobResult(jobResult)
		}
		scheduleAfter = scheduler.TrySchedule()
		scheduleTimer.Reset(scheduleAfter)
	}
}

func (scheduler *Scheduler) PushJobEvent(jobEvent *common.JobEvent)  {
	scheduler.JobEventChan <- jobEvent
}

func (scheduler *Scheduler) PushJobResult(result *common.JobExecuteResult) {
	var (
		err        error
		redisKey   string
	)

	scheduler.JobExecuteChan <- result
	redisKey = common.REDIS_CRON_RESULT + result.ExecuteInfo.Job.Name
	if _, err = library.GRedisServer.SetCacheData(redisKey, result, 3600*24); err != nil {
		fmt.Println("set job result fail")
	}
}

func InitScheduler() (err error) {
	GScheduler = &Scheduler{
		JobEventChan:      make(chan *common.JobEvent, 1000),
		JobPlanTable:      make(map[string]*common.JobSchedulerPlan),
		JobExecutingTable: make(map[string]*common.JobExecuteInfo),
		JobExecuteChan:    make(chan *common.JobExecuteResult, 1000),
	}

	go GScheduler.schedulerLoop()
	return
}