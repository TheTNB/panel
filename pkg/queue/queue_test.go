package queue

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
)

type QueueTestSuite struct {
	suite.Suite
}

func TestQueueTestSuite(t *testing.T) {
	suite.Run(t, &QueueTestSuite{})
}

func (suite *QueueTestSuite) TestQueueInitialization() {
	queue := New(10)
	suite.NotNil(queue)
	suite.NotNil(queue.jobs)
}

func (suite *QueueTestSuite) TestPushJobToQueue() {
	queue := New(10)
	job := &MockJob{}
	err := queue.Push(job, []any{"arg1", "arg2"})
	suite.NoError(err)
}

func (suite *QueueTestSuite) TestBulkJobsToQueue() {
	queue := New(10)
	jobs := []JobItem{
		{Job: &MockJob{}, Args: []any{"arg1"}},
		{Job: &MockJob{}, Args: []any{"arg2"}},
	}
	err := queue.Bulk(jobs)
	suite.NoError(err)
}

func (suite *QueueTestSuite) TestLaterJobExecution() {
	queue := New(10)
	job := &MockJob{}
	err := queue.Later(1, job, []any{"arg1"})
	suite.NoError(err)
}

func (suite *QueueTestSuite) TestRunQueue() {
	queue := New(10)
	job := &MockJob{}
	suite.NoError(queue.Push(job, []any{"arg1"}))
	queue.Run(context.Background())
	time.Sleep(1 * time.Second)
	suite.True(job.Executed)
}

func (suite *QueueTestSuite) TestRunQueueWithLaterJob() {
	queue := New(10)
	job := &MockJob{}
	suite.NoError(queue.Later(1, job, []any{"arg1"}))
	queue.Run(context.Background())
	time.Sleep(2 * time.Second)
	suite.True(job.Executed)
}

func (suite *QueueTestSuite) TestRunQueueWithBulkJobs() {
	queue := New(10)
	jobs := []JobItem{
		{Job: &MockJob{}, Args: []any{"arg1"}},
		{Job: &MockJob{}, Args: []any{"arg2"}},
	}
	suite.NoError(queue.Bulk(jobs))
	queue.Run(context.Background())
	time.Sleep(1 * time.Second)
}

func (suite *QueueTestSuite) TestRunQueueWithErrHandle() {
	queue := New(10)
	job := &MockJob{}
	suite.NoError(queue.Push(job, []any{"arg1"}))
	queue.Run(context.Background())
	time.Sleep(1 * time.Second)
	suite.Error(job.Err)
}

type MockJob struct {
	Executed bool
	Err      error
}

func (job *MockJob) Handle(args ...any) error {
	job.Executed = true
	return errors.New("error")
}

func (job *MockJob) ErrHandle(err error) {
	job.Err = err
}
