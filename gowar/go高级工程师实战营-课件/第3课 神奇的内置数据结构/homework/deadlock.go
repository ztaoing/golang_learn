package main

import "sync"

type task struct{}

type MyMap struct {
	m   map[int]task
	mux sync.RWMutex
}

func (m *MyMap) finishJob(t task, id int) {
	m.mux.Lock()
	defer m.mux.Unlock()

	// finish task
	delete(m.m, id)
}

func (m *MyMap) DoMyJob(taskID int) {
	m.mux.RLock()
	defer m.mux.RUnlock()

	t := m.m[taskID]

	m.finishJob(t, taskID)
}

func main() {
	var taskMap = &MyMap{
		m: map[int]task{},
	}
	taskMap.DoMyJob(1)
}
