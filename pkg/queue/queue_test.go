package queue

import (
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
	queue := New()
	suite.NotNil(queue)
	suite.NotNil(queue.jobs)
	suite.NotNil(queue.isShutdown)
	suite.NotNil(queue.done)
}

func (suite *QueueTestSuite) TestPushJobToQueue() {
	queue := New()
	job := &MockJob{}
	err := queue.Push(job, []any{"arg1", "arg2"})
	suite.NoError(err)
}

func (suite *QueueTestSuite) TestPushJobToShutdownQueue() {
	queue := New()
	queue.Run()
	suite.NoError(queue.Shutdown())
	job := &MockJob{}
	err := queue.Push(job, []any{"arg1", "arg2"})
	suite.Error(err)
	suite.EqualError(err, "queue is shutdown, cannot add new jobs")
}

func (suite *QueueTestSuite) TestBulkJobsToQueue() {
	queue := New()
	jobs := []Jobs{
		{Job: &MockJob{}, Args: []any{"arg1"}},
		{Job: &MockJob{}, Args: []any{"arg2"}},
	}
	err := queue.Bulk(jobs)
	suite.NoError(err)
}

func (suite *QueueTestSuite) TestBulkJobsToShutdownQueue() {
	queue := New()
	queue.Run()
	suite.NoError(queue.Shutdown())
	jobs := []Jobs{
		{Job: &MockJob{}, Args: []any{"arg1"}},
		{Job: &MockJob{}, Args: []any{"arg2"}},
	}
	err := queue.Bulk(jobs)
	suite.Error(err)
	suite.EqualError(err, "queue is shutdown, cannot add new jobs")
}

func (suite *QueueTestSuite) TestLaterJobExecution() {
	queue := New()
	job := &MockJob{}
	err := queue.Later(1, job, []any{"arg1"})
	suite.NoError(err)
}

func (suite *QueueTestSuite) TestLaterJobExecutionOnShutdownQueue() {
	queue := New()
	queue.Run()
	suite.NoError(queue.Shutdown())
	job := &MockJob{}
	err := queue.Later(1, job, []any{"arg1"})
	suite.NoError(err)
}

func (suite *QueueTestSuite) TestRunQueue() {
	queue := New()
	job := &MockJob{}
	suite.NoError(queue.Push(job, []any{"arg1"}))
	queue.Run()
	time.Sleep(1 * time.Second)
	suite.True(job.Executed)
}

func (suite *QueueTestSuite) TestRunQueueWithLaterJob() {
	queue := New()
	job := &MockJob{}
	suite.NoError(queue.Later(1, job, []any{"arg1"}))
	queue.Run()
	time.Sleep(2 * time.Second)
	suite.True(job.Executed)
}

func (suite *QueueTestSuite) TestRunQueueWithBulkJobs() {
	queue := New()
	jobs := []Jobs{
		{Job: &MockJob{}, Args: []any{"arg1"}},
		{Job: &MockJob{}, Args: []any{"arg2"}},
	}
	suite.NoError(queue.Bulk(jobs))
	queue.Run()
	time.Sleep(1 * time.Second)
}

func (suite *QueueTestSuite) TestRunQueueWithErrHandle() {
	queue := New()
	job := &MockJob{}
	suite.NoError(queue.Push(job, []any{"arg1"}))
	queue.Run()
	time.Sleep(1 * time.Second)
	suite.Error(job.Err)
}

func (suite *QueueTestSuite) TestShutdownQueue() {
	queue := New()
	queue.Run()
	err := queue.Shutdown()
	suite.NoError(err)
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
