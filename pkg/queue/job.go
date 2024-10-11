package queue

type Job interface {
	Handle(args ...any) error
}

type JobWithErrHandle interface {
	Job
	ErrHandle(err error)
}

type JobItem struct {
	Job   Job
	Args  []any
	Delay uint
}
