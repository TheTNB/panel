package bootstrap

import (
	"github.com/TheTNB/panel/pkg/queue"
)

func NewQueue() *queue.Queue {
	return queue.New(100)
}
