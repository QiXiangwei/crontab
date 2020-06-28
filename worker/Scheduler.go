package worker

import (
	"crontab/common"
	"fmt"
	"time"
)

type Scheduler struct {
	JobEventChan chan *common.JobEvent
	JobPlanTable map[string]*common.JobSchedulerPlan
}

var (
	GScheduler *Scheduler
)

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
			fmt.Println("go cron:" + jobPlan.Job.Name)
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

func (scheduler *Scheduler) schedulerLoop() {
	var (
		jobEvent *common.JobEvent
		scheduleAfter time.Duration
		scheduleTimer *time.Timer
	)

	scheduleAfter = scheduler.TrySchedule()
	scheduleTimer = time.NewTimer(scheduleAfter)

	for {
		select {
		case jobEvent = <- scheduler.JobEventChan:
		scheduler.handleJobEvent(jobEvent)
		case <- scheduleTimer.C:
		}
		scheduleAfter = scheduler.TrySchedule()
		scheduleTimer.Reset(scheduleAfter)
	}
}

func (scheduler *Scheduler) PushJobEvent(jobEvent *common.JobEvent)  {
	scheduler.JobEventChan <- jobEvent
}

func InitScheduler() (err error) {
	GScheduler = &Scheduler{
		JobEventChan: make(chan *common.JobEvent),
		JobPlanTable: make(map[string]*common.JobSchedulerPlan),
	}



	go GScheduler.schedulerLoop()
	return
}