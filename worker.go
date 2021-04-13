package scheduler

import (
	"time"
)

var (
	workers = map[int64]*Worker{}
)

type Worker struct {
	WorkerID int64
	Tasks    map[string]*Task
	Timer    *time.Timer
	Schedule
}

func RunSchedule(task *Task) *Worker {
	workerID := task.RunAt.Unix()
	if _, exists := workers[workerID]; !exists {
		workers[workerID] = &Worker{
			WorkerID: workerID,
			Tasks:    make(map[string]*Task),
		}
	}

	// Add task to Worker
	workers[workerID].Tasks[task.TaskID] = task

	if workers[workerID].Timer == nil {
		workers[workerID].Timer = time.NewTimer(task.RunAt.Sub(time.Now()))
		go workers[workerID].Start()
	}

	return workers[workerID]
}

func (w *Worker) Start() {
	<-w.Timer.C

	for _, task := range w.Tasks {
		go task.Run()
	}
}
