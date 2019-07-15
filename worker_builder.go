package lucy

import "golang.org/x/xerrors"

// WorkerBuilder is the builder of Worker.
type WorkerBuilder struct {
	WorkerQueue         WorkerQueue
	Logger              Logger
	HTTPClient          HTTPClient
	RequestMiddlewares  []func(request *Request)
	ResponseMiddlewares []func(response *Response)
	Spider              func(response *Response) ([]*Request, error)
}

// NewWorkerBuilder is builder of the WorkerBuilder that initialize fields by default values.
func NewWorkerBuilder() *WorkerBuilder {
	workerBuilder := &WorkerBuilder{}
	return workerBuilder
}

// Build builds worker from given fields.
func (w *WorkerBuilder) Build() (*Worker, error) {
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
		w.RequestMiddlewares = make([]func(*Request), 0)
	}
	if w.ResponseMiddlewares == nil {
		w.ResponseMiddlewares = make([]func(*Response), 0)
	}
	return newWorker(
		w.WorkerQueue,
		w.HTTPClient,
		w.Logger,
		w.RequestMiddlewares,
		w.ResponseMiddlewares,
		w.Spider,
	), nil
}

// SetWorkerQueue sets WorkerQueue implementation
func (w *WorkerBuilder) SetWorkerQueue(workerQueue WorkerQueue) {
	w.WorkerQueue = workerQueue
}

// SetLogger sets Logger imeplementation
func (w *WorkerBuilder) SetLogger(logger Logger) {
	w.Logger = logger
}

// SetHTTPClient sets HTTPClient implementation
func (w *WorkerBuilder) SetHTTPClient(httpClient HTTPClient) {
	w.HTTPClient = httpClient
}

// AddRequestMiddleware adds requestMiddleware in requestMiddleware list
func (w *WorkerBuilder) AddRequestMiddleware(f func(request *Request)) {
	if w.RequestMiddlewares == nil {
		w.RequestMiddlewares = make([]func(request *Request), 0)
	}
	w.RequestMiddlewares = append(w.RequestMiddlewares, f)
}

// AddResponseMiddleware adds responseMiddleware in responseMiddleware list
func (w *WorkerBuilder) AddResponseMiddleware(f func(response *Response)) {
	if w.ResponseMiddlewares == nil {
		w.ResponseMiddlewares = make([]func(response *Response), 0)
	}
	w.ResponseMiddlewares = append(w.ResponseMiddlewares, f)
}

// SetSpider sets spider
func (w *WorkerBuilder) SetSpider(f func(response *Response) ([]*Request, error)) {
	w.Spider = f
}
