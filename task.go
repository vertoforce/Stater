package stater

import (
	"context"
	"fmt"
)

// Errors
var (
	ErrTaskDone     = fmt.Errorf("task already done")
	ErrTaskNotFound = fmt.Errorf("task not found")
)

// Task is some long running task that can be started, stopped, and resumed
//
// Note that the entire task structure must be serializable
type Task struct {
	// Some way to uniquely identify this task
	ID string

	// Current State
	State         State
	workFunction  IncrementalWorkFunction
	storageDriver StorageDriver
	messager      *Messager

	// Current running work function context
	ctx    context.Context
	cancel context.CancelFunc

	// Task state
	Done    bool
	Running bool
}

// IncrementalWorkFunction performs the smallest possible portion of work, returning the new state
//
// If the work function returns an error, but the state != nil, the new state will be used on the next call.
// If the state is nil though, it will exist the task and mark it done.
//
// If the returned state is nil, the original state will be used on the next call
type IncrementalWorkFunction func(ctx context.Context, state State, messager *Messager) (State, error)

// NewTask Creates a new task that has some starting state, and an incredmental function
func NewTask(ID string, startingState State, workFunction IncrementalWorkFunction, messager *Messager) *Task {
	return &Task{
		ID:           ID,
		State:        startingState,
		workFunction: workFunction,
		messager:     messager,
	}
}

// Start will start performing the task as fast as possible.
// It will keep calling the incremental work function.
// It will not return until the task is done
// After each run of the incremental function, it will store the state using the storage driver
// Once the work function returned state and error are nil, the task will be marked done.
//
// If the task is started once it is already running, it will do nothing.
// If the task is done, it will return an error
func (t *Task) Start(ctx context.Context, storageDriver StorageDriver) error {
	if t.Running {
		return nil
	}
	if t.Done {
		return ErrTaskDone
	}
	t.storageDriver = storageDriver
	t.Running = true

	for t.Running {
		workCtx, cancel := context.WithCancel(ctx)
		t.ctx = workCtx
		t.cancel = cancel
		newState, err := t.workFunction(workCtx, t.State, t.messager)
		cancel()
		if err != nil {
			if !t.Running {
				// We were paused immediately, not truly errored
				return nil
			}
			t.markDone()
			return err
		}
		if newState == nil {
			// Task is done
			t.markDone()
			break
		}
		t.updateState(newState)
	}

	return nil
}

func (t *Task) updateState(state State) {
	t.storageDriver.SaveTask(t)
}

// markDone is called when the task is deemed done
func (t *Task) markDone() {
	t.messager.SendMessage(&Message{
		Type:    DoneMessage,
		Task:    t,
		Message: "We finished",
	})
	t.Running = false
	t.Done = true
	t.storageDriver.RemoveTask(t.ID)
}

// Pause the running task.  It will not resume until Resume() is called
//
// If the task is already paused this will do nothing
func (t *Task) Pause() {
	t.Running = false
}

// PauseImmediately cancels the current working context to stop the work as quickly as possible.
//
// Note that this will likely leave a IncrementalWorkFunction partly done
func (t *Task) PauseImmediately() {
	t.Pause()
	t.cancel()
}
