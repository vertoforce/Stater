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
	tasks, messager, err := stater.NewTaskEngine(context.Background(), fileDriver)
	if err != nil {
		t.Error(err)
		return
	}

	// Create thread to listen for message
	gotAMessage := false
	go func() {
		messageStream := messager.GetMessageStream()
		for {
			message := <-messageStream
			if message.Task.ID == "MyTask" {
				gotAMessage = true
			}
		}
	}()

	// Check to make sure we don't have any tasks
	if len(tasks) > 0 {
		t.Errorf("Too many tasks in storage")
		return
	}

	workFunction := func(ctx context.Context, state stater.State, messager *stater.Messager) (stater.State, error) {
		currentCount := state.GetState()["Count"].(int)
		if currentCount == 5 {
			return nil, nil
		}
		newCount := currentCount + 1
		newState := state
		newState.GetState()["Count"] = newCount
		return newState, nil
	}
	task := stater.NewTask("MyTask", &stater.BasicState{Fields: map[string]interface{}{"Count": 0}}, workFunction, messager)
	task.Start(context.Background(), fileDriver)

	// Wait for task to finish
	time.Sleep(time.Millisecond * 100)

	// Check that the task did finish
	if !gotAMessage {
		t.Errorf("Task did not send completion message")
	}
}
