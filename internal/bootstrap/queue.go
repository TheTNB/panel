package bootstrap

import (
	"context"

	"github.com/TheTNB/panel/internal/app"
	"github.com/TheTNB/panel/pkg/queue"
)

func initQueue() {
	app.Queue = queue.New(40)
	go app.Queue.Run(context.Background())
}
