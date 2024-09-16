package queue

import (
	"errors"
	"time"
)

type Queue struct {
	jobs       chan Jobs
	isShutdown chan struct{}
	done       chan struct{}
}

func NewQueue(bufferSize ...int) *Queue {
	if len(bufferSize) == 0 {
		bufferSize = append(bufferSize, 100)
	}
	return &Queue{
		jobs:       make(chan Jobs, bufferSize[0]),
		isShutdown: make(chan struct{}),
		done:       make(chan struct{}),
	}
}

func (r *Queue) Push(job Job, args []any) error {
	select {
	case <-r.isShutdown:
		return errors.New("queue is shutdown, cannot add new jobs")
	default:
		r.jobs <- Jobs{Job: job, Args: args}
		return nil
	}
}

func (r *Queue) Bulk(jobs []Jobs) error {
	for _, job := range jobs {
		if job.Delay > 0 {
			time.AfterFunc(time.Duration(job.Delay)*time.Second, func() {
				select {
				case <-r.isShutdown:
					return
				default:
					r.jobs <- Jobs{Job: job.Job, Args: job.Args}
				}
			})
			continue
		}

		select {
		case <-r.isShutdown:
			return errors.New("queue is shutdown, cannot add new jobs")
		default:
			r.jobs <- job
		}
	}

	return nil
}

func (r *Queue) Later(delay uint, job Job, args []any) error {
	time.AfterFunc(time.Duration(delay)*time.Second, func() {
		select {
		case <-r.isShutdown:
			return
		default:
			r.jobs <- Jobs{Job: job, Args: args}
		}
	})

	return nil
}

func (r *Queue) Run() {
	go func() {
		for {
			select {
			case job := <-r.jobs:
				if err := job.Job.Handle(job.Args...); err != nil {
					if errJob, ok := job.Job.(JobWithErrHandle); ok {
						errJob.ErrHandle(err)
					}
				}
			case <-r.isShutdown:
				close(r.done)
				return
			}
		}
	}()
}

func (r *Queue) Shutdown() error {
	close(r.isShutdown)
	<-r.done
	return nil
}
