package queue

import (
	"context"
	"sync"

	"github.com/getumen/lucy"
	"github.com/wangjia184/sortedset"
)

var cond *sync.Cond

func init() {
	cond = sync.NewCond(&sync.Mutex{})
}

type memoryWorkerQueue struct {
	queue *sortedset.SortedSet
}

// NewMemoryWorkerQueue return Memory WorkerQueue implementation
func NewMemoryWorkerQueue() (lucy.WorkerQueue, error) {
	return &memoryWorkerQueue{
		queue: sortedset.New(),
	}, nil
}

func (q *memoryWorkerQueue) SubscribeRequests(ctx context.Context) (<-chan *lucy.Request, error) {
	requestChan := make(chan *lucy.Request)

	go func() {
		defer close(requestChan)
		for {
			var ok bool
			select {
			case _, ok = <-ctx.Done():
			default:
				ok = true
			}
			if !ok {
				break
			}
			cond.L.Lock()
			node := q.queue.PopMin()
			if node != nil {
				request, ok := node.Value.(*lucy.Request)
				if ok {
					requestChan <- request
				}
			} else {
				cond.Wait()
			}
			cond.L.Unlock()
		}
	}()

	return requestChan, nil
}
func (q *memoryWorkerQueue) RetryRequest(request *lucy.Request) error {
	cond.L.Lock()
	defer cond.L.Unlock()
	if q.queue.GetByKey(request.URL) == nil {
		q.queue.AddOrUpdate(request.URL, sortedset.SCORE(request.Priority), request)
		cond.Signal()
	}
	return nil
}

func (q *memoryWorkerQueue) PublishRequest(request *lucy.Request) error {
	cond.L.Lock()
	defer cond.L.Unlock()
	if q.queue.GetByKey(request.URL) == nil {
		q.queue.AddOrUpdate(request.URL, sortedset.SCORE(request.Priority), request)
		cond.Signal()
	}
	return nil
}
