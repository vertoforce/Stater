# Stater

[![Go Report Card](https://goreportcard.com/badge/github.com/vertoforce/stater)](https://goreportcard.com/report/github.com/vertoforce/stater)
[![Documentation](https://godoc.org/github.com/vertoforce/stater?status.svg)](https://godoc.org/github.com/vertoforce/stater)

Stater is a library to help perform very long running tasks.  It supports reboots, pausing, and resuming.

The library works by allowing you to define a `IncrementalWorkFunction` that performs the smallest amount of work possible, and then updating the Task's _state_.  After each run of the `IncrementalWorkFunction` the task stores the state using your defined _StorageDriver_.

You can then pause and resume the task.  If the program restarts, on the next start the task engine will recognize the abandoned tasks, and start them again.

## Usage

The goal is to ultimately call

```go
task := TaskEngine.NewTask(...)
task.Start(...)
```

That function will start a task that will run even on restarts.

However to get there you must first

1. Create some storage driver (where we will store the state)
2. Define your work functions
3. Create a task engine

Each task is relatively painless, but is necessary to start the task off of the engine.

The [godoc](https://godoc.org/github.com/vertoforce/stater) has a full example of doing the above steps.

Below are some more instructions

### Create a storage driver

The storage driver is a way to store/load the state of running tasks.  You can create your own that obeys the [StorageDriver](https://godoc.org/github.com/vertoforce/Stater#StorageDriver) interface, or just use the built in `filedriver` (see example in godoc).

### Define your work functions

We define each work function as a map so that we can save the "name" of the work function for each task.  This is because after each work done by the task, we cannot save the work function to disk, we need to save a "reference/pointer" so we know what to run when the program starts again.

To define your work functions simple create a `map[string]stater.IncrementalWorkFunction`

### Create a task engine

The task engine is the central location to store the [StorageDriver](https://godoc.org/github.com/vertoforce/Stater#StorageDriver) and the [Messager](https://godoc.org/github.com/vertoforce/Stater#Messager).

It also starts all paused tasks if the program restarts.

To create it simply use the StorageDriver and workFunctions you have defined above.

## Important notes

This library is early in development and will have bugs.  I'm sure there are some interesting race conditions possible, but for simple usage it should hold up.
