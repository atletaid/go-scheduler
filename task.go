package scheduler

import (
	"fmt"
	"math/rand"
	"reflect"
	"time"
)

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const charsetLength = len(charset)

var seededRand *rand.Rand = rand.New(
	rand.NewSource(time.Now().UnixNano()))

type Schedule struct {
	Timer         *time.Timer
	RunAt         time.Time
	IsRunInterval bool
	RunEvery      time.Duration
	RunUntil      time.Time
}

type Task struct {
	TaskID string
	Schedule
	Func   FunctionMeta
	Params []Param
}

func NewTask(function FunctionMeta, params []Param) *Task {
	str := make([]byte, 4)
	for i := range str {
		str[i] = charset[seededRand.Intn(charsetLength)]
	}

	return &Task{
		TaskID: fmt.Sprint(string(str), time.Now().UnixNano()),
		Func:   function,
		Params: params,
	}
}

func (task *Task) SetTime(runAt time.Time) {
	task.RunAt = runAt
	task.Timer = time.NewTimer(runAt.Sub(time.Now()))
}

func (task *Task) SetInterval(runEvery time.Duration, runUntil time.Time) {
	task.RunAt = time.Now().Add(runEvery)
	task.Timer = time.NewTimer(task.RunAt.Sub(time.Now()))
	task.IsRunInterval = true
	task.RunEvery = runEvery
	task.RunUntil = runUntil
}

func (task *Task) Run() {
	<-task.Timer.C
	function := reflect.ValueOf(task.Func.function)
	params := make([]reflect.Value, len(task.Params))
	for i, param := range task.Params {
		params[i] = reflect.ValueOf(param)
	}
	function.Call(params)

	if task.IsRunInterval && (time.Now().Before(task.RunUntil.Add(time.Second))) {
		task.SetTime(time.Now().Add(task.RunEvery))
		task.Run()
	}
}

func (task *Task) Stop() {
	task.Timer.Stop()
}
