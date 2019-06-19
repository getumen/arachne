package lucy

import (
	"context"

	"golang.org/x/sync/semaphore"
	"golang.org/x/xerrors"
)

// Worker handles scheduled requests.
type Worker struct {
	workerQueue                WorkerQueue
	requestRestrictionStrategy RequestRestrictionStrategy
	requestSemaphore           RequestSemaphore
	logger                     Logger
	maxRequestNum              int64
}

func newWorker(workerQueue WorkerQueue) *Worker {
	return &Worker{
		workerQueue: workerQueue,
	}
}

// Start kicks worker off.
func (*Worker) Start(ctx context.Context) {

}

func (w *Worker) subscribe(ctx context.Context) (<-chan Request, error) {
	output := make(chan Request)
	requestChan, err := w.workerQueue.SubscribeRequests(ctx)
	if err != nil {
		return nil, xerrors.Errorf("fail to subscribe requests: %w", err)
	}
	go func() {
		defer close(output)

		for request := range requestChan {
			output <- request
		}
	}()
	return output, nil
}

func (w *Worker) doRequest(requestChan <-chan Request) (<-chan Response, error) {
	responseChan := make(chan Response)

	// TODO: add goroutine supervisor
	go func() {
		workerRequestSemaphore := semaphore.NewWeighted(w.maxRequestNum)
		for request := range requestChan {
			ctx := context.Background()
			err := workerRequestSemaphore.Acquire(ctx, 1)
			if err != nil {
				w.logger.Errorf("fail to acquire workerRequestSemaphore")
				continue
			}

			go func(request *Request) {
				defer workerRequestSemaphore.Release(1)

				// request restriction
				if w.requestRestrictionStrategy.CheckRestriction() {
					resource, err := w.requestRestrictionStrategy.Resource(request)
					if err != nil {
						w.logger.Warnf("fail to get resource name for semaphore. %v", err)
						return
					}
					ctx := context.TODO()
					err = w.requestSemaphore.Acquire(ctx, resource)
					if err != nil {
						w.logger.Infof("retry %s.", request.URL)
						w.requestRestrictionStrategy.ChangePriorityWhenRestricted(request)
						err = w.workerQueue.RetryRequest(request)
						w.logger.Errorf("fail to retry %s. this request is lost.")
						return
					}
					defer w.requestSemaphore.Release(resource)
				}

				// TODO: do request

				// TODO: parse response and send
			}(&request)
		}
	}()
	return responseChan, nil
}

func (w *Worker) applySpider(responseChan <-chan Response) <-chan Request {
	requestChan := make(chan Request)

	go func() {
		// TODO: apply spider
	}()

	return requestChan
}

func (w *Worker) publishRequest(requestChan <-chan Request) error {
	// TODO: publish request
	return nil
}

// WorkerBuilder is the builder of Worker.
type WorkerBuilder struct {
}
