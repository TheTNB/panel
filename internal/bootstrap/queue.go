package bootstrap

import (
	"github.com/tnb-labs/panel/pkg/queue"
)

func NewQueue() *queue.Queue {
	return queue.New(100)
}
