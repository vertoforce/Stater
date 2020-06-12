package stater_test

import (
	"context"
	"fmt"

	"github.com/vertoforce/stater"
	"github.com/vertoforce/stater/storageDrivers/filedriver"
)

func Example() {
	// Create storage driver
	fileDriver, _ := filedriver.NewFileDriver("state.json")

	// Create your list of worker function
	// We need to define this to the engine since we cannot serialize the function with the task
	// We can only serialize the "name" of the worker function
	workerFunctions := map[string]stater.IncrementalWorkFunction{
		"example": func(ctx context.Context, state *stater.State, messager *stater.Messager) (*stater.State, error) {
			// We use floats because numbers will always unmarshal to a float
			currentCount := state.Fields["Count"].(float64)
			if currentCount == 5 {
				// This job is done!
				return nil, nil
			}
			state.Fields["Count"] = currentCount + 1
			return state, nil
		},
	}

	// Create a task engine
	// The task engine will restart and previously running tasks that were stored in the storage engine
	taskEngine, _ := stater.NewTaskEngine(context.Background(), fileDriver, workerFunctions)

	// Start a new task on the engine
	task := taskEngine.NewTask("MyTask", &stater.State{Fields: map[string]interface{}{"Count": float64(0)}}, "example")

	// Optionally save this task right now (it will save after every incremental function)
	fileDriver.SaveTask(task)

	// Start a thread to listen to messages
	go func() {
		message := <-taskEngine.Messager.GetMessageStream()
		fmt.Printf("%s:%s\n", message.Task.ID, message.Type)
	}()

	// Start the task, and wait for it to finish
	// At this point, even if our program crashes, once it restarts the task will start working again
	task.Start(context.Background())

	// Output: MyTask:done
}
