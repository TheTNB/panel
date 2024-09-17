package middleware

import (
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
		panic(err)
	}

	limiter, err := httplimit.NewMiddleware(store, httplimit.IPKeyFunc())
	if err != nil {
		panic(err)
	}

	return limiter.Handle
}
