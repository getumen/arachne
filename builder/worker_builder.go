package builder

import (
	"github.com/getumen/lucy"
	"golang.org/x/xerrors"
)

// WorkerBuilder is the builder of Worker.
type WorkerBuilder struct {
	WorkerQueue         lucy.WorkerQueue
	Logger              lucy.Logger
	HTTPClient          lucy.HTTPClient
	RequestMiddlewares  []func(request *lucy.Request)
	ResponseMiddlewares []func(response *lucy.Response)
	Spider              func(response *lucy.Response) ([]*lucy.Request, error)
}

// NewWorkerBuilder is builder of the WorkerBuilder that initialize fields by default values.
func NewWorkerBuilder() *WorkerBuilder {
	workerBuilder := &WorkerBuilder{}
	return workerBuilder
}

// Build builds worker from given fields.
func (w *WorkerBuilder) Build() (*lucy.Worker, error) {
	if w.WorkerQueue == nil {
		return nil, xerrors.New("set worker queue")
	}
	if w.Logger == nil {
		return nil, xerrors.New("set logger")
	}
	if w.HTTPClient == nil {
		return nil, xerrors.New("set http client")
	}
	if w.Spider == nil {
		return nil, xerrors.New("set spider")
	}
	if w.RequestMiddlewares == nil {
		w.RequestMiddlewares = make([]func(*lucy.Request), 0)
	}
	if w.ResponseMiddlewares == nil {
		w.ResponseMiddlewares = make([]func(*lucy.Response), 0)
	}
	return &lucy.Worker{
		WorkerQueue:         w.WorkerQueue,
		HTTPClient:          w.HTTPClient,
		Logger:              w.Logger,
		RequestMiddlewares:  w.RequestMiddlewares,
		ResponseMiddlewares: w.ResponseMiddlewares,
		Spider:              w.Spider,
	}, nil
}

// SetWorkerQueue sets WorkerQueue implementation
func (w *WorkerBuilder) SetWorkerQueue(workerQueue lucy.WorkerQueue) *WorkerBuilder {
	w.WorkerQueue = workerQueue
	return w
}

// SetLogger sets Logger imeplementation
func (w *WorkerBuilder) SetLogger(logger lucy.Logger) *WorkerBuilder {
	w.Logger = logger
	return w
}

// SetHTTPClient sets HTTPClient implementation
func (w *WorkerBuilder) SetHTTPClient(httpClient lucy.HTTPClient) *WorkerBuilder {
	w.HTTPClient = httpClient
	return w
}

// SetSpider sets spider
func (w *WorkerBuilder) SetSpider(f func(response *lucy.Response) ([]*lucy.Request, error)) *WorkerBuilder {
	w.Spider = f
	return w
}
