package stater

import "context"

// StorageDriver is some storage mechanism that can store and retrieve tasks
type StorageDriver interface {
	// Save or update a task
	SaveTask(task *Task) error
	RemoveTask(ID string) error
	LoadTasks() ([]*Task, error)
}

// TaskEngine is a collection of tasks, worker functions, and messager
type TaskEngine struct {
	Messager        *Messager
	StorageDriver   StorageDriver
	Tasks           []*Task
	workerFunctions map[string]IncrementalWorkFunction
}

// NewTaskEngine Gets all the stored tasks from a storage driver and starts them in a new thread
//
// It will create a messager to use for future created tasks
//
// This function should always be called first
func NewTaskEngine(ctx context.Context, storageDriver StorageDriver, workerFunctions map[string]IncrementalWorkFunction) (taskEngine *TaskEngine, err error) {
	tasks, err := storageDriver.LoadTasks()
	if err != nil {
		return nil, err
	}
	messager := NewMessager()

	engine := &TaskEngine{
		workerFunctions: workerFunctions,
		Tasks:           tasks,
		Messager:        messager,
	}

	// Start each task
	for _, task := range tasks {
		// Set connections
		task.engine = engine
		task.messager = messager

		go func(ctx context.Context, task *Task) {
			task.Start(ctx, storageDriver)
		}(ctx, task)
	}

	return engine, nil
}
