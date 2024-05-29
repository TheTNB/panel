package internal

type Task interface {
	Process(taskID uint) error
	DispatchWaiting() error
}
