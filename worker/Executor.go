package worker

import (
	"context"
	"crontab/common"
	"crontab/library"
	"os/exec"
	"time"
)

type Executor struct {

}

var (
	GExecutor *Executor
)

func (executor *Executor) ExecuteJob(info *common.JobExecuteInfo) {
	go func() {
		var (
			cmd           *exec.Cmd
			output        []byte
			err           error
			executeResult *common.JobExecuteResult
			jobLock       *library.Lock
		)

		executeResult = &common.JobExecuteResult{
			ExecuteInfo: info,
			StartTime:   time.Now(),
		}

		jobLock = library.GEtcServer.CreateEtcLock(info.Job.Name)

		executeResult.StartTime = time.Now()
		err = jobLock.TryLock()
		defer jobLock.UnLock()
		if err != nil {
			executeResult.Err      = err
			executeResult.StopTime = time.Now()
		} else {
			cmd         = exec.CommandContext(context.TODO(), "/bin/bash", "-c", info.Job.Command)
			output, err = cmd.CombinedOutput()

			executeResult.Output   = output
			executeResult.Err      = err
			executeResult.StopTime = time.Now()

			GScheduler.PushJobResult(executeResult)
		}
	}()
}

func InitExecutor() (err error) {
	GExecutor = &Executor{}
	return
}
