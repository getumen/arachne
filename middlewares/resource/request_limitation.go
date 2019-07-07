package resource

import (
	"context"

	"github.com/getumen/lucy"
	"golang.org/x/sync/semaphore"
)

// RequestCounter is a request counter
type RequestCounter struct {
	sema *semaphore.Weighted
}

// NewRequestCounter creates new request counter
func NewRequestCounter(maxRequest int64) *RequestCounter {
	return &RequestCounter{
		sema: semaphore.NewWeighted(maxRequest),
	}
}

// RequestMiddleware is a request middleware
func (r *RequestCounter) RequestMiddleware(*lucy.Request) {
	ctx := context.TODO()
	r.sema.Acquire(ctx, 1)
}

// ResponseMiddleware is a response middleware
func (r *RequestCounter) ResponseMiddleware(*lucy.Response) {
	r.sema.Release(1)
}
