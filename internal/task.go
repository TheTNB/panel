package internal

type Task interface {
	Process(taskID uint)
	DispatchWaiting() error
}
