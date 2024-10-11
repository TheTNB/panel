package queue

import (
	"context"
	"errors"
	"time"
)

type Queue struct {
	jobs chan JobItem
}

func New(bufferSize int) *Queue {
	return &Queue{
		jobs: make(chan JobItem, bufferSize),
	}
}

func (r *Queue) Push(job Job, args []any) error {
	select {
	case r.jobs <- JobItem{Job: job, Args: args}:
		return nil
	default:
		return errors.New("job queue is full")
	}
}

func (r *Queue) Bulk(jobs []JobItem) error {
	for _, job := range jobs {
		jobCopy := job
		if jobCopy.Delay > 0 {
			time.AfterFunc(time.Duration(jobCopy.Delay)*time.Second, func() {
				r.jobs <- jobCopy
			})
			continue
		}

		select {
		case r.jobs <- jobCopy:
			return nil
		default:
			return errors.New("job queue is full")
		}
	}

	return nil
}

func (r *Queue) Later(delay uint, job Job, args []any) error {
	jobCopy := job
	argsCopy := make([]any, len(args))
	copy(argsCopy, args)
	time.AfterFunc(time.Duration(delay)*time.Second, func() {
		r.jobs <- JobItem{Job: jobCopy, Args: argsCopy}
	})

	return nil
}

func (r *Queue) Run(ctx context.Context) {
	go func() {
		for {
			select {
			case job := <-r.jobs:
				processJob(job)
			case <-ctx.Done():
				return
			}
		}
	}()
}

func (r *Queue) Len() int {
	return len(r.jobs)
}

func (r *Queue) IsFull() bool {
	return len(r.jobs) == cap(r.jobs)
}

func processJob(job JobItem) {
	if err := job.Job.Handle(job.Args...); err != nil {
		if errJob, ok := job.Job.(JobWithErrHandle); ok {
			errJob.ErrHandle(err)
		}
	}
}
