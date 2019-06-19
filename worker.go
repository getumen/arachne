package lucy

import (
	"context"
	"golang.org/x/xerrors"
)

type Worker struct {
	workerQueue WorkerQueue
}

func newWorker(workerQueue WorkerQueue) *Worker {
	return &Worker{
		workerQueue: workerQueue,
	}
}

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

		for {
			select {
			case <-ctx.Done():
				return
			case request := <-requestChan:
				output <- request
			}
		}
	}()
	return output, nil
}

func (w *Worker) doRequest() {

}

type WorkerBuilder struct {
}
