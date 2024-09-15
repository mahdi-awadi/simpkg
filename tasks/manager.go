package tasks

import (
	"errors"
	"sort"
)

// IManager interface
type IManager interface {
	IsStarted() bool
	Current() *Task
	Items() Items
	Add(*Task)
	Get(string) *Task
	Once(string) error
	SetCurrent(*Task)
	OnStatusChange(func(*Task))
	OnStart(func())
	OnStop(func())
	OnComplete(func())
	Start()
	Stop()
	Reset()
}

// Manager struct
type Manager struct {
	onTaskStatusChange []func(*Task)

	items      Items
	onStart    func()
	onStop     func()
	onComplete func()
	current    *Task
	isSorted   bool
	isStarted  bool
}

// Items sortable
type Items []*Task

// New returns task manager new instance
func New() IManager {
	return &Manager{
		items:              make(Items, 0),
		onTaskStatusChange: make([]func(*Task), 0),
	}
}

// IsStarted returns start status
func (manager *Manager) IsStarted() bool {
	return manager.isStarted
}

// Current current task
func (manager *Manager) Current() *Task {
	return manager.current
}

// Items return task items
func (manager *Manager) Items() Items {
	return manager.items
}

// OnStatusChange set on status change
func (manager *Manager) OnStatusChange(fn func(*Task)) {
	manager.onTaskStatusChange = append(manager.onTaskStatusChange, fn)
}

// OnStart trigger on task starts
func (manager *Manager) OnStart(fn func()) {
	manager.onStart = fn
}

// OnStop trigger on task stops
func (manager *Manager) OnStop(fn func()) {
	manager.onStop = fn
}

// OnComplete trigger on task complete
func (manager *Manager) OnComplete(fn func()) {
	manager.onComplete = fn
}

// triggerOnTaskStatusChange trigger on task status change
func (manager *Manager) triggerOnTaskStatusChange(task *Task) {
	if len(manager.onTaskStatusChange) > 0 {
		for _, handler := range manager.onTaskStatusChange {
			handler(task)
		}
	}
}

// Add task
func (manager *Manager) Add(task *Task) {
	task.status = StatusPending
	taskOnstart := task.OnStart
	task.OnStart = func() {
		manager.triggerOnTaskStatusChange(task)
		if taskOnstart != nil {
			taskOnstart()
		}
	}

	taskOnRetry := task.OnRetry
	task.OnRetry = func(err error, index int) {
		manager.triggerOnTaskStatusChange(task)
		if taskOnRetry != nil {
			taskOnRetry(err, index)
		}
	}

	// task run end
	task.onRunEnd = func(task *Task) {
		manager.triggerOnTaskStatusChange(task)
		if task.err != nil {
			manager.isStarted = false
			return
		}

		if !task.executeNextTask {
			manager.isStarted = false
			return
		}

		// run next task
		manager.runNext(task)
	}

	// append
	manager.items = append(manager.items, task)
}

// runNext run next task
func (manager *Manager) runNext(current *Task) {
	if current != nil && current.IsError() {
		manager.isStarted = false
		return
	}

	// if stopped
	if !manager.isStarted {
		return
	}

	// find next task
	currentIndex := -1
	if current != nil {
		currentIndex = current.index
	}

	task, completed := manager.move(true, currentIndex)
	if completed {
		manager.isStarted = false
		if manager.onComplete != nil {
			manager.onComplete()
		}
		return
	}

	// if task already completed
	if task == nil || (!completed && task == nil) {
		manager.isStarted = false
		return
	}

	task.reset()
	if current != nil {
		current.reset()
	}

	manager.SetCurrent(task)
	task.Start()
}

// move to the next/prev task
func (manager *Manager) move(next bool, current int) (*Task, bool) {
	if next {
		current++
	} else {
		current--
		if current < 0 {
			current = 0
		}
	}

	if current > len(manager.items)-1 {
		return nil, true
	}

	for _, task := range manager.items {
		if task.index == current {
			return task, false
		}
	}

	return nil, false
}

// SetCurrent set current task
func (manager *Manager) SetCurrent(task *Task) {
	manager.current = task
}

// Get find specified task
func (manager *Manager) Get(name string) *Task {
	for _, task := range manager.items {
		if task.Name == name {
			return task
		}
	}

	return nil
}

// Once run task once
func (manager *Manager) Once(name string) error {
	task := manager.Get(name)
	if task == nil {
		return errors.New("task not found")
	}

	task.reset()
	task.executeNextTask = false
	task.Start()
	task.executeNextTask = true
	return task.err
}

// Start executes the tasks
func (manager *Manager) Start() {
	manager.sort()
	if len(manager.items) == 0 {
		return
	}

	manager.isStarted = true
	if manager.onStart != nil {
		manager.onStart()
	}

	if manager.current == nil {
		manager.runNext(nil)
		return
	}

	// run current task
	manager.current.Start()
}

// Stop stops the tasks
func (manager *Manager) Stop() {
	manager.isStarted = false
	if manager.onStop != nil {
		manager.onStop()
	}
}

// Reset resets the manager
func (manager *Manager) Reset() {
	manager.isStarted = false
	manager.current = nil
	for _, task := range manager.items {
		task.reset()
	}
}

// sort TaskItems by order
func (manager *Manager) sort() {
	if !manager.isSorted {
		sort.Sort(manager.items)
		for index, task := range manager.items {
			task.index = index
		}
		manager.isSorted = true
	}
}

func (v Items) Len() int           { return len(v) }
func (v Items) Swap(i, j int)      { v[i], v[j] = v[j], v[i] }
func (v Items) Less(i, j int) bool { return v[i].Order < v[j].Order }
