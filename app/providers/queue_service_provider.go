package providers

import (
	"github.com/goravel/framework/contracts/foundation"
	"github.com/goravel/framework/contracts/queue"
	"github.com/goravel/framework/facades"

	"github.com/TheTNB/panel/v2/app/jobs"
)

type QueueServiceProvider struct {
}

func (receiver *QueueServiceProvider) Register(app foundation.Application) {
	if err := facades.Queue().Register(receiver.Jobs()); err != nil {
		panic(err.Error())
	}
}

func (receiver *QueueServiceProvider) Boot(app foundation.Application) {

}

func (receiver *QueueServiceProvider) Jobs() []queue.Job {
	return []queue.Job{
		&jobs.ProcessTask{},
	}
}
