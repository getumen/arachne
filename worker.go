package lucy

import (
	"context"
	"golang.org/x/xerrors"
)

type Worker struct {
	workerQueue WorkerQueue
}

type WorkerBuilder struct {
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

	SubscribeLoop:
		for {
			select {
			case <-ctx.Done():
				break SubscribeLoop
			case request := <-requestChan:
				output <- request
			}
		}
	}()
	return output, nil
}

func (w *Worker) doRequest() {

}
