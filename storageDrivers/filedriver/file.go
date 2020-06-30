package filedriver

import (
	"os"
)

// FileDriver stores the tasks in a local file
type FileDriver struct {
	File *os.File
}

// NewFileDriver Creates a new file driver using the provided filename
func NewFileDriver(fileName string) (*FileDriver, error) {
	file, err := os.OpenFile(fileName, os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		return nil, err
	}
	return NewFileDriverFromFile(file)
}

// NewFileDriverFromFile Creates a new file driver using the provided file
func NewFileDriverFromFile(file *os.File) (*FileDriver, error) {
	driver := &FileDriver{
		File: file,
	}
	// If there are no tasks, make an empty file
	if tasks, err := driver.LoadTasks(); err != nil || len(tasks) == 0 {
		file.Seek(0, 0)
		file.Truncate(0)
		_, err := file.Write([]byte("[]"))
		if err != nil {
			return nil, err
		}
	}

	return driver, nil
}
