package builder

import (
	"github.com/getumen/arachne"
	"golang.org/x/xerrors"
)

// WorkerBuilder is the builder of Worker.
type WorkerBuilder struct {
	WorkerQueue         arachne.WorkerQueue
	Logger              arachne.Logger
	HTTPClient          arachne.HTTPClient
	RequestMiddlewares  []func(request *arachne.Request)
	ResponseMiddlewares []func(response *arachne.Response)
	Spider              func(response *arachne.Response) ([]*arachne.Request, error)
}

// NewWorkerBuilder is builder of the WorkerBuilder that initialize fields by default values.
func NewWorkerBuilder() *WorkerBuilder {
	workerBuilder := &WorkerBuilder{}
	return workerBuilder
}

// Build builds worker from given fields.
func (w *WorkerBuilder) Build() (*arachne.Worker, error) {
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
		w.RequestMiddlewares = make([]func(*arachne.Request), 0)
	}
	if w.ResponseMiddlewares == nil {
		w.ResponseMiddlewares = make([]func(*arachne.Response), 0)
	}
	return &arachne.Worker{
		WorkerQueue:         w.WorkerQueue,
		HTTPClient:          w.HTTPClient,
		Logger:              w.Logger,
		RequestMiddlewares:  w.RequestMiddlewares,
		ResponseMiddlewares: w.ResponseMiddlewares,
		Spider:              w.Spider,
	}, nil
}

// SetWorkerQueue sets WorkerQueue implementation
func (w *WorkerBuilder) SetWorkerQueue(workerQueue arachne.WorkerQueue) *WorkerBuilder {
	w.WorkerQueue = workerQueue
	return w
}

// SetLogger sets Logger imeplementation
func (w *WorkerBuilder) SetLogger(logger arachne.Logger) *WorkerBuilder {
	w.Logger = logger
	return w
}

// SetHTTPClient sets HTTPClient implementation
func (w *WorkerBuilder) SetHTTPClient(httpClient arachne.HTTPClient) *WorkerBuilder {
	w.HTTPClient = httpClient
	return w
}

// SetSpider sets spider
func (w *WorkerBuilder) SetSpider(f func(response *arachne.Response) ([]*arachne.Request, error)) *WorkerBuilder {
	w.Spider = f
	return w
}
