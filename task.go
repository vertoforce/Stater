package stater

import (
	"context"
	"fmt"
	"sync"
)

// Errors
var (
	ErrTaskDone     = fmt.Errorf("task already done")
	ErrTaskNotFound = fmt.Errorf("task not found")
)

// Task is some long running task that can be started, stopped, and resumed
//
// Note that the entire task structure must be serializable, including the state
type Task struct {
	// Some way to uniquely identify this task
	ID string

	// Current State
	State            *State
	WorkFunctionName string
	lock             sync.Mutex
	engine           *TaskEngine

	// Current running work function context
	ctx    context.Context
	cancel context.CancelFunc

	// Task state
	Done    bool
	Running bool
}

// IncrementalWorkFunction performs the smallest possible portion of work, returning the new state
//
// If the work function returns an error, the task will finish and return the error
// If the state is nil, the task will be done
//
// If the returned state is nil, the original state will be used on the next call
type IncrementalWorkFunction func(ctx context.Context, state *State, messager *Messager) (*State, error)

// NewTask Creates a new task that has some starting state, and an incredmental function
func (e *TaskEngine) NewTask(ID string, startingState *State, workFunctionName string) *Task {
	t := &Task{
		ID:               ID,
		State:            startingState,
		WorkFunctionName: workFunctionName,
		lock:             sync.Mutex{},
		engine:           e,
	}

	return t
}

// Start will start performing the task as fast as possible.
// It will keep calling the incremental work function.
// It will not return until the task is done.
// After each run of the incremental function, it will store the state using the storage driver.
// Once the work function returned state and error are nil, the task will be marked done.
//
// If the task is started once it is already running, it will do nothing.
// If the task is done, it will return an error.
func (t *Task) Start(ctx context.Context) error {
	t.lock.Lock()
	if t.Running {
		t.lock.Unlock()
		return nil
	}
	if t.Done {
		t.lock.Unlock()
		return ErrTaskDone
	}
	t.Running = true
	t.lock.Unlock()
	workFunction, ok := t.engine.workerFunctions[t.WorkFunctionName]
	if !ok {
		t.markDone()
		return fmt.Errorf("worker function not found")
	}

	for t.Running {
		workCtx, cancel := context.WithCancel(ctx)
		t.ctx = workCtx
		t.cancel = cancel
		newState, err := workFunction(workCtx, t.State, t.engine.Messager)
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

// updateState stores the new state of the task
func (t *Task) updateState(state *State) {
	t.engine.StorageDriver.SaveTask(t)
}

// markDone is called when the task is deemed done
//
// It will set running=false, done=true, remove the task, and send a DoneMessage
// to the engine messager
func (t *Task) markDone() {
	t.engine.Messager.SendMessage(&Message{
		Type:    DoneMessage,
		Task:    t,
		Message: "We finished",
	})
	t.lock.Lock()
	t.Running = false
	t.Done = true
	t.engine.StorageDriver.RemoveTask(t.ID)
	t.lock.Unlock()
}

// Pause the running task.  It will not resume until Resume() is called
//
// If the task is already paused this will do nothing
func (t *Task) Pause() {
	t.lock.Lock()
	t.Running = false
	t.lock.Unlock()
}

// PauseImmediately cancels the current working context to stop the work as quickly as possible.
//
// Note that this will likely leave a IncrementalWorkFunction partly done
func (t *Task) PauseImmediately() {
	t.lock.Lock()
	t.Pause()
	t.lock.Unlock()
	t.cancel()
}
