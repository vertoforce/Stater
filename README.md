# Stater

Stater is a library to help perform very long running tasks.  It supports reboots, pausing, and resuming.

The library works by allowing you to define a `IncrementalWorkFunction` that performs the smallest amount of work possible, and then updating the Task's _state_.  After each run of the `IncrementalWorkFunction` the task stores the state of the task (using whatever storage driver you want).

You can then pause and resume the task.  If the program restart, on the next start the task engine will recognize there are stored tasks, and start them again.

## Usage

See the godoc for examples.  I did my best to documents things as I understand some things may be confusing.  If you open an issue and I'd be glad to help.
