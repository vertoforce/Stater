package stater

import "context"

// StorageDriver is some storage mechanism that can store and retrieve tasks
type StorageDriver interface {
	// Save or update a task
	SaveTask(task *Task) error
	RemoveTask(ID string) error
	LoadTasks() ([]*Task, error)
}

// NewTaskEngine Gets all the stored tasks from a storage driver and starts them in a new thread
//
// It will create a messager to use for future created tasks
//
// This function should always be called first
func NewTaskEngine(ctx context.Context, storageDriver StorageDriver) (tasks []*Task, messager *Messager, err error) {
	tasks, err = storageDriver.LoadTasks()
	if err != nil {
		return nil, nil, err
	}
	messager = NewMessager()

	// Start each task
	for _, task := range tasks {
		go func(ctx context.Context, task *Task) {
			task.Start(ctx, storageDriver)
		}(ctx, task)
	}

	return tasks, messager, nil
}
