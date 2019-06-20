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
	requestMiddlewares         []func(request *Request)
	responseMidlewares         []func(response *Response)
	httpClient                 HTTPClient
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
						w.logger.Warnf("fail to get resource name for semaphore.: %v", err)
						return
					}
					ctx := context.TODO()
					err = w.requestSemaphore.Acquire(ctx, resource)
					if err != nil {
						w.logger.Infof("retry %s because worker failed to acquire resource.", request.URL)
						w.requestRestrictionStrategy.ChangePriorityWhenRestricted(request)
						err = w.workerQueue.RetryRequest(request)
						w.logger.Errorf("fail to retry %s. this request is lost.")
						return
					}
					defer w.requestSemaphore.Release(resource)
				}

				// apply requestMiddlewares
				for _, middlewareFunc := range w.requestMiddlewares {
					middlewareFunc(request)
					if request == nil {
						// discard request if nil
						return
					}
				}

				// send request
				httpRequest, err := request.HTTPRequest()
				if err != nil {
					w.logger.Warnf("fail to construct http.Request. %v: %v", request, err)
					return
				}
				httpResponse, err := w.httpClient.Do(httpRequest)
				if err != nil {
					w.logger.Warnf("fail to get http.Response of http.Request(%v): %v", request, err)
					return
				}
				response, err := NewResponseFromHTTPResponse(httpResponse)
				if err != nil {
					w.logger.Warnf("fail to construct Response of http.Response(%v): %v", httpResponse, err)
					return
				}

				// apply responseMiddlewares
				for _, middlewareFunc := range w.responseMidlewares {
					middlewareFunc(response)
					if request == nil {
						// discard response if nil
						return
					}
				}

				responseChan <- *response
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
