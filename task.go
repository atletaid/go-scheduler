package scheduler

import (
	"fmt"
	"math/rand"
	"reflect"
	"time"
)

const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
const charsetLength = len(charset)

var seededRand *rand.Rand = rand.New(
	rand.NewSource(time.Now().UnixNano()))

type Schedule struct {
	RunAt    time.Time
	RunEvery time.Duration
	RunUntil time.Time
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

func (task *Task) SetRunAt(runAt time.Time) *Task {
	task.Schedule = Schedule{
		RunAt: runAt,
	}
	return task
}

func (task *Task) SetInterval(runAt time.Time, runEvery time.Duration, runUntil time.Time) *Task {
	task.Schedule = Schedule{
		RunAt:    runAt,
		RunEvery: runEvery,
		RunUntil: runUntil,
	}
	return task
}

func (task *Task) Run() {
	function := reflect.ValueOf(task.Func.function)
	params := make([]reflect.Value, len(task.Params))
	for i, param := range task.Params {
		params[i] = reflect.ValueOf(param)
	}
	function.Call(params)

	if time.Now().Before(task.RunUntil.Add(time.Second)) {
		task.RunAt = time.Now().Add(task.RunEvery)
		RunSchedule(task)
	}
}
