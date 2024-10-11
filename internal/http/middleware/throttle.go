package middleware

import (
	"log"
	"net/http"
	"time"

	"github.com/sethvargo/go-limiter/httplimit"
	"github.com/sethvargo/go-limiter/memorystore"
)

// Throttle 限流器
func Throttle(tokens uint64, interval time.Duration) func(next http.Handler) http.Handler {
	store, err := memorystore.New(&memorystore.Config{
		Tokens:   tokens,
		Interval: interval,
	})
	if err != nil {
		log.Fatalf("failed to create throttle memorystore: %v", err)
	}

	limiter, err := httplimit.NewMiddleware(store, httplimit.IPKeyFunc())
	if err != nil {
		log.Fatalf("failed to initialize throttle middleware: %v", err)
	}

	return limiter.Handle
}
