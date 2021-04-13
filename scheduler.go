package scheduler

import (
	"fmt"
	"time"
)

type Scheduler struct {
	funcRegistry *FuncRegistry
}

func NewScheduler() *Scheduler {
	return &Scheduler{
		funcRegistry: NewFuncRegistry(),
	}
}

func (scheduler *Scheduler) RunAt(runAt time.Time, function Function, params ...Param) (string, error) {
	funcMeta, err := scheduler.funcRegistry.Add(function)
	if err != nil {
		return "", err
	}

	task := NewTask(funcMeta, params).SetRunAt(runAt)
	RunSchedule(task)

	return task.TaskID, nil
}

func (scheduler *Scheduler) RunAfter(duration time.Duration, function Function, params ...Param) (string, error) {
	return scheduler.RunAt(time.Now().Add(duration), function, params...)
}

func (scheduler *Scheduler) RunEvery(runEvery time.Duration, runUntil time.Time, function Function, params ...Param) (string, error) {
	funcMeta, err := scheduler.funcRegistry.Add(function)
	if err != nil {
		return "", err
	}

	task := NewTask(funcMeta, params).SetInterval(time.Now(), runEvery, runUntil)
	RunSchedule(task)

	return task.TaskID, nil
}

func (scheduler *Scheduler) GetAllScheduler() {
	return
}

func (scheduler *Scheduler) Stop(taskID string) error {
	for _, worker := range workers {
		if _, exists := worker.Tasks[taskID]; exists {
			delete(worker.Tasks, taskID)

			if len(worker.Tasks) == 0 {
				worker.Timer.Stop()
			}

			return nil
		}
	}

	return fmt.Errorf("Task %v not found", taskID)
}

func (scheduler *Scheduler) ClearAll() {
	// for _, task := range scheduler.tasks {
	// 	task.Stop()
	// 	delete(scheduler.tasks, task.TaskID)
	// }
	// scheduler.funcRegistry = NewFuncRegistry()
}

func (scheduler *Scheduler) Reschedule(taskID string, time time.Time) error {
	// task, found := scheduler.tasks[taskID]
	// if !found {
	// 	return fmt.Errorf("Task %v not found", taskID)
	// }

	// task.Stop()
	// task.SetNextRun(time)

	// go task.Run()
	return nil
}
