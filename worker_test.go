package lucy

import (
	"context"
	"github.com/golang/mock/gomock"
	"golang.org/x/xerrors"
	"testing"
)

func setupTestWorker_subscribe(t *testing.T) (chan Request, func()) {
	retChan := make(chan Request)

	go func() {
		for _, urlString := range []string{
			"https://golang.org/",
			"https://golang.org/doc/",
			"https://golang.org/pkg/",
		} {
			retChan <- *NewGetRequest(urlString)
		}
	}()

	return retChan, func() {
		close(retChan)
	}
}

func setupTestWorker_subscribeInfiniteChannel(t *testing.T) (chan Request, func()) {
	retChan := make(chan Request)

	running := true

	go func() {
		for running {
			retChan <- *NewGetRequest("https://golang.org/")
		}
	}()

	return retChan, func() {
		running = false
		close(retChan)
	}
}

func TestWorker_subscribeSuccess(t *testing.T) {

	subscribeChan, tearDown := setupTestWorker_subscribe(t)
	defer tearDown()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()

	workerQueueMock := NewMockWorkerQueue(ctrl)
	workerQueueMock.EXPECT().SubscribeRequests(ctx).Return(subscribeChan, nil)

	worker := newWorker(workerQueueMock)
	reqChan, err := worker.subscribe(ctx)

	if err != nil {
		t.Fatalf("subscribe fail %v", err)
	}

	tests := []string{
		"https://golang.org/",
		"https://golang.org/doc/",
		"https://golang.org/pkg/",
	}

	for i := 0; i < len(tests); i++ {
		actual := <-reqChan
		if actual.URL != tests[i] {
			t.Fatalf("expected %s, but got %s", tests[i], actual.URL)
		}
	}
}

func TestWorker_subscribeError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()

	workerQueueMock := NewMockWorkerQueue(ctrl)
	workerQueueMock.EXPECT().SubscribeRequests(ctx).Return(nil, xerrors.New("some error."))

	worker := newWorker(workerQueueMock)
	reqChan, err := worker.subscribe(ctx)

	if reqChan != nil || err == nil {
		t.Fatalf("fail to propagate error.")
	}
}
