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
func (f *FileDriver) RemoveTask(t *stater.Task) error {
	// TODO:
	return nil
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

	// Re-write tasks
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
