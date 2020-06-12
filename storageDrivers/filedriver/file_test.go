package filedriver

import (
	"os"
	"testing"

	"github.com/vertoforce/stater"
)

func TestFile(t *testing.T) {
	file, err := os.OpenFile("state.json", os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0644)
	if err != nil {
		t.Error(err)
		return
	}

	driver, err := NewFileDriverFromFile(file)
	if err != nil {
		t.Error(err)
		return
	}

	tasks, err := driver.LoadTasks()
	if err != nil {
		t.Error(err)
		return
	}
	if len(tasks) > 0 {
		t.Errorf("Was not expecting tasks")
	}

	// Save a task
	err = driver.SaveTask(stater.NewTask("Test", nil, nil, nil))
	if err != nil {
		t.Error(err)
		return
	}

	// Load tasks
	tasks, err = driver.LoadTasks()
	if err != nil {
		t.Error(err)
		return
	}
	if len(tasks) != 1 {
		t.Errorf("Task did not save")
		return
	}
	if tasks[0].ID != "Test" {
		t.Errorf("Task did not save correctly")
	}

	// Remove the task
	err = driver.RemoveTask(tasks[0].ID)
	if err != nil {
		t.Error(err)
		return
	}

	tasks, err = driver.LoadTasks()
	if err != nil {
		t.Error(err)
		return
	}
	if len(tasks) > 0 {
		t.Errorf("Was not expecting tasks")
	}
}
