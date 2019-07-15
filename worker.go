package lucy

import (
	"context"
	"sync"

	"golang.org/x/xerrors"
)

// channelSize is worker chennel size
const channelSize = 16

// Worker handles scheduled requests.
type Worker struct {
	workerQueue         WorkerQueue
	httpClient          HTTPClient
	logger              Logger
	requestMiddlewares  []func(request *Request)
	responseMiddlewares []func(response *Response)
	spider              func(response *Response) ([]*Request, error)
}

func newWorker(
	workerQueue WorkerQueue,
	httpClient HTTPClient,
	logger Logger,
	requestMiddlewares []func(request *Request),
	responseMiddlewares []func(response *Response),
	spider func(response *Response) ([]*Request, error),
) *Worker {
	return &Worker{
		workerQueue:         workerQueue,
		httpClient:          httpClient,
		logger:              logger,
		requestMiddlewares:  requestMiddlewares,
		responseMiddlewares: responseMiddlewares,
		spider:              spider,
	}
}

// Start kicks worker off.
func (w *Worker) Start(ctx context.Context) error {
	requestPipeline, err := w.subscribe(ctx)
	if err != nil {
		return xerrors.Errorf("fail to subscribe request: %v", err)
	}
	responsePipeline, err := w.doRequest(requestPipeline)
	if err != nil {
		return xerrors.Errorf("fail to initialize request pipeline: %v", err)
	}
	nextRequestPipeline, err := w.applySpider(responsePipeline)
	if err != nil {
		return xerrors.Errorf("fail to initialize spider pipeline: %v", err)
	}
	err = w.publishRequest(nextRequestPipeline)
	return xerrors.Errorf("fail to publish request: %v", err)
}

// StartWithFirstRequest kicks worker off with first request.
func (w *Worker) StartWithFirstRequest(ctx context.Context, URL string) error {
	requestPipeline, err := w.subscribe(ctx)
	if err != nil {
		return xerrors.Errorf("fail to subscribe request: %v", err)
	}
	responsePipeline, err := w.doRequest(requestPipeline)
	if err != nil {
		return xerrors.Errorf("fail to initialize request pipeline: %v", err)
	}
	nextRequestPipeline, err := w.applySpider(responsePipeline)
	if err != nil {
		return xerrors.Errorf("fail to initialize spider pipeline: %v", err)
	}

	request, err := NewGetRequest(URL)
	if err != nil {
		return xerrors.Errorf("fail to create initial request: %v", err)
	}
	nextRequestPipeline <- request

	err = w.publishRequest(nextRequestPipeline)
	return xerrors.Errorf("fail to publish request: %v", err)
}

func (w *Worker) subscribe(ctx context.Context) (<-chan *Request, error) {
	output := make(chan *Request, channelSize)
	requestChan, err := w.workerQueue.SubscribeRequests(ctx)
	if err != nil {
		return nil, xerrors.Errorf("fail to subscribe requests: %w",
			err)
	}
	go func() {
		defer close(output)

		for request := range requestChan {
			w.logger.Debugf("subscribe %s", request.URL)
			output <- request
		}
	}()
	return output, nil
}

func (w *Worker) doRequest(requestChan <-chan *Request) (<-chan *Response, error) {
	responseChan := make(chan *Response, channelSize)

	// TODO: add goroutine supervisor
	go func() {
		defer close(responseChan)

		requestWaitGroup := sync.WaitGroup{}

		for request := range requestChan {

			handleRequest := func(request *Request) {
				defer requestWaitGroup.Done()

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
					w.logger.Warnf("fail to construct http.Request. %v: %v",
						request, err)
					return
				}
				w.logger.Debugf("request %s", httpRequest.URL.String())
				httpResponse, err := w.httpClient.Do(httpRequest)
				if err != nil {
					w.logger.Warnf("fail to get http.Response of http.Request(%v): %v",
						request, err)
					return
				}
				response, err := NewResponseFromHTTPResponse(httpResponse)
				if err != nil {
					w.logger.Warnf("fail to construct Response of http.Response(%v): %v",
						httpResponse, err)
					return
				}

				// apply responseMiddlewares
				for _, middlewareFunc := range w.responseMiddlewares {
					middlewareFunc(response)
					if response == nil {
						// discard response if nil
						return
					}
				}

				responseChan <- response
			}

			requestWaitGroup.Add(1)
			go handleRequest(request)
		}

		requestWaitGroup.Wait()
	}()

	return responseChan, nil
}

func (w *Worker) applySpider(responseChan <-chan *Response) (chan *Request, error) {
	requestChan := make(chan *Request, channelSize)

	go func() {
		defer close(requestChan)
		for response := range responseChan {
			w.logger.Debugf("apply Spider to %s", response.Request.URL)
			nextRequests, err := w.spider(response)
			if err != nil {
				w.logger.Infof("spider error: %v", err)
				continue
			}
			for _, request := range nextRequests {
				requestChan <- request
			}
		}
	}()

	return requestChan, nil
}

func (w *Worker) publishRequest(requestChan <-chan *Request) error {
	for request := range requestChan {
		w.logger.Debugf("publish %s", request.URL)
		err := w.workerQueue.PublishRequest(request)
		if err != nil {
			w.logger.Errorf("fail to publish request: %s", request.URL)
		}
	}
	return nil
}

// RetryMiddleware is request middleware that remove request in worker pipeline
// and send request to worker queue if Request.Meta['retry'] flas is true.
func (w *Worker) RetryMiddleware(request *Request) {
	if retry, ok := request.Meta["retry"]; ok {
		if retryFlag, ok := retry.(bool); retryFlag && ok {
			err := w.workerQueue.RetryRequest(request)
			if err != nil {
				w.logger.Errorf("fail to retry %s. this request is lost.")
			}
			request = nil
		}
	}
}
