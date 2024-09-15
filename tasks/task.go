package tasks

import "time"

// Status is Task status
type Status string

const (
	StatusPending Status = "pending"
	StatusStart   Status = "start"
	StatusRetry   Status = "retry"
	StatusFail    Status = "fail"
	StatusSuccess Status = "success"
)

// ITask interface
type ITask interface {
	Status() Status
	Error() error
	IsDone() bool
	IsSuccess() bool
	IsError() bool
	Start()
	Attempts() int
}

// Task struct
type Task struct {
	index   int
	Order   int
	Name    string
	Label   string
	Handler func() error

	RetryDelay     time.Duration
	RetryCondition func(err error, index int) bool
	Retries        int

	Message        any
	StartMessage   func(*Task) any
	RetryMessage   func(*Task, error, int) any
	ErrorMessage   func(*Task, error) any
	SuccessMessage func(*Task) any

	OnStart   func()
	OnRetry   func(err error, index int)
	OnError   func(error)
	OnSuccess func()
	onRunEnd  func(*Task)

	err             error
	status          Status
	executeNextTask bool
	attempts        int
}

// Status returns task status
func (task *Task) Status() Status {
	return task.status
}

// SetStatus set status
func (task *Task) SetStatus(status Status) {
	task.status = status
}

// Error returns task err
func (task *Task) Error() error {
	return task.err
}

// Attempts returns task attempts
func (task *Task) Attempts() int {
	return task.attempts
}

// IsError returns true if task has error
func (task *Task) IsError() bool {
	return task.err != nil
}

// IsDone returns true if task is done
func (task *Task) IsDone() bool {
	return task.status == StatusSuccess || task.status == StatusFail
}

// IsSuccess returns true if task is success
func (task *Task) IsSuccess() bool {
	return task.status == StatusSuccess
}

// pause to next retry
func (task *Task) pause() {
	if task.RetryDelay > 0 {
		time.Sleep(task.RetryDelay)
	}
}

// reset statues
func (task *Task) reset() {
	task.err = nil
	task.SetStatus(StatusPending)
	task.executeNextTask = true
	task.attempts = 0
}

// Start the task
func (task *Task) Start() {
	onTaskEnd := func() {
		if task.onRunEnd != nil {
			task.onRunEnd(task)
		}
	}

	// if task already done
	if task.IsSuccess() {
		onTaskEnd()
		return
	}

	task.attempts = 0
	if task.OnStart != nil {
		task.SetStatus(StatusStart)
		if task.StartMessage != nil {
			task.Message = task.StartMessage(task)
		}
		task.OnStart()
	}

	if task.Retries == 0 {
		task.Retries = 1
	}

	for index := 0; index < task.Retries; index++ {
		task.attempts++
		if index > 0 {
			task.SetStatus(StatusRetry)
			if task.RetryMessage != nil {
				task.Message = task.RetryMessage(task, task.err, index)
			}

			// retry condition
			if task.RetryCondition != nil && task.err != nil && !task.RetryCondition(task.err, index) {
				break
			}

			// pause
			task.pause()
			if task.OnRetry != nil {
				task.OnRetry(task.err, index)
			}
		}

		// run the task and get result
		result := task.Handler()
		err, hasErr := result.(error)
		if !hasErr {
			task.err = nil
			break
		}

		task.err = err
		task.SetStatus(StatusFail)
		if task.ErrorMessage != nil {
			task.Message = task.ErrorMessage(task, err)
		}
	}

	// on error
	if task.err != nil && task.OnError != nil {
		task.OnError(task.err)
	}

	// on success
	if task.err == nil && task.OnSuccess != nil {
		task.OnSuccess()
	}

	if task.err == nil {
		task.SetStatus(StatusSuccess)
		if task.SuccessMessage != nil {
			task.Message = task.SuccessMessage(task)
		}
	} else {
		task.SetStatus(StatusFail)
		if task.ErrorMessage != nil {
			task.Message = task.ErrorMessage(task, task.err)
		}
	}

	// at run end
	onTaskEnd()
}
