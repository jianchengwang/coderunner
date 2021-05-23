package task

import "time"

type Sandbox struct {
	T *Task
	LastOptTime time.Time
}

var SandboxMap = make(map[string]Sandbox)