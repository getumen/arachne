package lucy

import (
	"net/http"

	"github.com/getumen/lucy"
	"github.com/getumen/lucy/logger"
	"github.com/getumen/lucy/queue"
	"golang.org/x/xerrors"
)

// WorkerBuilder is the builder of Worker.
type WorkerBuilder struct {
	workerQueue         lucy.WorkerQueue
	logger              lucy.Logger
	httpClient          lucy.HTTPClient
	requestMiddlewares  []func(request *lucy.Request)
	responseMiddlewares []func(response *lucy.Response)
	spider              func(response *lucy.Response) ([]*lucy.Request, error)
}

// NewWorkerBuilder is builder of the WorkerBuilder that initialize fields by default values.
func NewWorkerBuilder() (*WorkerBuilder, error) {
	workerBuilder := &WorkerBuilder{}
	inMemoryQueue, err := queue.NewMemoryWorkerQueue()
	if err != nil {
		return nil, xerrors.Errorf("fail to create default queue: %w", err)
	}
	workerBuilder.workerQueue = inMemoryQueue
	logger := logger.NewStdoutLogger(lucy.InfoLevel)
	workerBuilder.logger = logger
	workerBuilder.httpClient = &http.Client{}
	workerBuilder.requestMiddlewares = make([]func(request *lucy.Request), 0)
	workerBuilder.responseMiddlewares = make([]func(response *lucy.Response), 0)

	return workerBuilder, nil
}
