package queue

import (
	"context"
	"sync"

	"github.com/getumen/lucy"
	"github.com/wangjia184/sortedset"
)

var memoryQueueMutex sync.RWMutex

func init() {
	memoryQueueMutex = sync.RWMutex{}
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
			memoryQueueMutex.RLock()
			node := q.queue.PopMin()
			if node != nil {
				request, ok := node.Value.(*lucy.Request)
				if ok {
					requestChan <- request
				}
			}
			memoryQueueMutex.RUnlock()
		}
	}()

	return requestChan, nil
}
func (q *memoryWorkerQueue) RetryRequest(request *lucy.Request) error {
	memoryQueueMutex.Lock()
	defer memoryQueueMutex.Unlock()
	if q.queue.GetByKey(request.URL) == nil {
		q.queue.AddOrUpdate(request.URL, sortedset.SCORE(request.Priority), request)
	}
	return nil
}

func (q *memoryWorkerQueue) PublishRequest(request *lucy.Request) error {
	memoryQueueMutex.Lock()
	defer memoryQueueMutex.Unlock()
	if q.queue.GetByKey(request.URL) == nil {
		q.queue.AddOrUpdate(request.URL, sortedset.SCORE(request.Priority), request)
	}
	return nil
}
