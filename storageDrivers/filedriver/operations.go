package filedriver

import (
	"encoding/json"

	"github.com/vertoforce/stater"
)

// LoadTasks from filedriver
func (f *FileDriver) LoadTasks() ([]*stater.Task, error) {
	_, err := f.File.Seek(0, 0)
	if err != nil {
		return nil, err
	}
	tasks := []*stater.Task{}

	// Load file tasks
	err = json.NewDecoder(f.File).Decode(&tasks)
	if err != nil {
		return nil, err
	}

	return tasks, nil
}

// RemoveTask from file driver
func (f *FileDriver) RemoveTask(ID string) error {
	tasks, err := f.LoadTasks()
	if err != nil {
		return err
	}

	// Find if this tasks exists
	for i, t := range tasks {
		if t.ID == ID {
			// Remove this task form the list
			tasks = append(tasks[0:i], tasks[i:len(tasks)-1]...)
		}
	}

	// Overwrite file
	return f.writeTasks(tasks)
}

// SaveTask to driver
func (f *FileDriver) SaveTask(task *stater.Task) error {
	tasks, err := f.LoadTasks()
	if err != nil {
		return err
	}

	// Find if this tasks exists
	updated := false
	for i, t := range tasks {
		if t.ID == task.ID {
			// We have this task, just overwrite it
			tasks[i] = task
			updated = true
		}
	}

	// We need to create this task
	if !updated {
		tasks = append(tasks, task)
	}

	return f.writeTasks(tasks)
}

func (f *FileDriver) writeTasks(tasks []*stater.Task) error {
	// Re-write tasks
	err := f.File.Truncate(0)
	if err != nil {
		return err
	}
	_, err = f.File.Seek(0, 0)
	if err != nil {
		return err
	}
	err = json.NewEncoder(f.File).Encode(tasks)
	if err != nil {
		return err
	}
	return nil
}

// LoadTask from file
func (f *FileDriver) LoadTask(ID string) (*stater.Task, error) {
	tasks, err := f.LoadTasks()
	if err != nil {
		return nil, err
	}

	// Find if this tasks exists
	for _, t := range tasks {
		if t.ID == ID {
			return t, nil
		}
	}
	return nil, stater.ErrTaskNotFound
}
