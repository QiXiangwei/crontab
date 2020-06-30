package worker

import (
	"context"
	"crontab/common"
	"os/exec"
	"time"
)

type Executor struct {

}

var (
	GExecutor *Executor
)

func (executor *Executor) ExecuteJob(info *common.JobExecuteInfo) {
	var (
		cmd           *exec.Cmd
		output        []byte
		err           error
		executeResult *common.JobExecuteResult
	)

	executeResult = &common.JobExecuteResult{
		ExecuteInfo: info,
		StartTime:   time.Now(),
	}
	cmd         = exec.CommandContext(context.TODO(), "/bin/bash", "-c", info.Job.Command)
	output, err = cmd.CombinedOutput()

	executeResult.Output   = output
	executeResult.Err      = err
	executeResult.StopTime = time.Now()

	GScheduler.PushJobResult(executeResult)
}

func InitExecutor() (err error) {
	GExecutor = &Executor{}
	return
}
