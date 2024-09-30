package bootstrap

import (
	"github.com/TheTNB/panel/internal/app"
	"github.com/TheTNB/panel/pkg/queue"
)

func initQueue() {
	app.Queue = queue.New()
	go app.Queue.Run()
}
