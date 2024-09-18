package bootstrap

import (
	"github.com/TheTNB/panel/internal/panel"
	"github.com/TheTNB/panel/pkg/queue"
)

func initQueue() {
	panel.Queue = queue.New()
	go panel.Queue.Run()
}
