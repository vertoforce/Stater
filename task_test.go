package stater_test

import (
	"context"
	"testing"
	"time"

	"github.com/vertoforce/stater"
	"github.com/vertoforce/stater/storageDrivers/filedriver"
)

func TestTask(t *testing.T) {
	// Create storage driver
	fileDriver, err := filedriver.NewFileDriver("state.json")
	if err != nil {
		t.Error(err)
		return
	}

	// Create a new task engine
	workFunction := func(ctx context.Context, state *stater.State, messager *stater.Messager) (*stater.State, error) {
		currentCount := state.Fields["Count"].(float64)
		if currentCount == 5 {
			return nil, nil
		}
		state.Fields["Count"] = currentCount + 1
		return state, nil
	}
	workerFunctions := map[string]stater.IncrementalWorkFunction{
		"one": workFunction,
	}
	taskEngine, err := stater.NewTaskEngine(context.Background(), fileDriver, workerFunctions)
	if err != nil {
		t.Error(err)
		return
	}

	// Create thread to listen for message
	gotAMessage := false
	go func() {
		messageStream := taskEngine.Messager.GetMessageStream()
		for {
			message := <-messageStream
			if message.Task.ID == "MyTask" {
				gotAMessage = true
			}
		}
	}()

	// Check to make sure we don't have any tasks
	tasks, err := taskEngine.StorageDriver.LoadTasks()
	if err != nil {
		t.Error(err)
		return
	}
	if len(tasks) > 0 {
		t.Errorf("Too many tasks in storage")
		return
	}

	task := taskEngine.NewTask("MyTask", &stater.State{Fields: map[string]interface{}{"Count": float64(0)}}, "one")
	// Start task and wait for it to finish
	task.Start(context.Background())

	// Check that the task did finish
	if !gotAMessage {
		t.Errorf("Task did not send completion message")
	}

	// Now try starting a task that was not done processing
	task = taskEngine.NewTask("MyTask2", &stater.State{Fields: map[string]interface{}{"Count": 0}}, "one")
	// Do not start the task (it'll finish too fast), instead mark it as "running" so it's as if it died when running
	task.Running = true
	err = fileDriver.SaveTask(task)
	if err != nil {
		t.Error(err)
		return
	}

	// Now we reload everything from the storage engine in to a new task engine
	// This will check if that task will start running again
	taskEngine, err = stater.NewTaskEngine(context.Background(), fileDriver, workerFunctions)
	if err != nil {
		t.Error(err)
		return
	}
	// Create thread to listen for message
	gotAMessage = false
	go func() {
		messageStream := taskEngine.Messager.GetMessageStream()
		for {
			message := <-messageStream
			if message.Task.ID == "MyTask2" {
				gotAMessage = true
			}
		}
	}()

	// Wait for task to re-start and complete
	time.Sleep(time.Millisecond * 500)

	if !gotAMessage {
		t.Error("Task did not re-start and complete")
	}

}
