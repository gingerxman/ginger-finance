package cron

import (
	"github.com/gingerxman/eel"
	"github.com/gingerxman/eel/cron"
)

type demoTask struct {
	cron.Task
}


func (this *demoTask) Run(taskCtx *cron.TaskContext) error {
	eel.Logger.Info("[demo_task] run...")
	return nil
}

func NewDemoTask() *demoTask{
	task := new(demoTask)
	task.Task = cron.NewTask("demo_task")
	return task
}

func init() {
	//task := NewDemoTask()
	//cron.RegisterTask(task, "*/5 * * * * *")
}