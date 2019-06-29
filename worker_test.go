package lucy

import (
	context "context"
	"errors"
	"net/http"
	"testing"

	gomock "github.com/golang/mock/gomock"
	"golang.org/x/xerrors"
)

func setupTestWorkerSubscribe(t *testing.T) (chan *Request, func()) {
	retChan := make(chan *Request)

	go func() {
		for _, urlString := range []string{
			"https://golang.org/",
			"https://golang.org/doc/",
			"https://golang.org/pkg/",
		} {
			r, _ := NewGetRequest(urlString)
			retChan <- r
		}
	}()

	return retChan, func() {
		close(retChan)
	}
}

func setupTestWorkerSubscribeInfiniteChannel(t *testing.T) (chan *Request, func()) {
	retChan := make(chan *Request)

	running := true

	go func() {
		for running {
			r, _ := NewGetRequest("https://golang.org/")
			retChan <- r
		}
	}()

	return retChan, func() {
		running = false
		close(retChan)
	}
}

func TestWorker_subscribeSuccess(t *testing.T) {

	subscribeChan, tearDown := setupTestWorkerSubscribe(t)
	defer tearDown()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()

	workerQueueMock := NewMockWorkerQueue(ctrl)
	workerQueueMock.EXPECT().SubscribeRequests(ctx).Return(subscribeChan, nil)

	worker := newWorker(workerQueueMock, nil, nil, nil, StdoutLogger{},
		10, []func(request *Request){}, []func(response *Response){})
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

	worker := newWorker(workerQueueMock, nil, nil, nil, StdoutLogger{},
		10, []func(request *Request){}, []func(response *Response){})
	reqChan, err := worker.subscribe(ctx)

	if reqChan != nil || err == nil {
		t.Fatalf("fail to propagate error.")
	}
}

// TestWorker_doRequestSuccess
// CheckRestriction: true
// semaphore resource: url domain
func TestWorker_doRequestRestrictionByDomain(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()

	requestStrategyMock := NewMockRequestRestrictionStrategy(ctrl)
	requestSemaphoreMock := NewMockRequestSemaphore(ctrl)
	httpClientMock := NewMockHTTPClient(ctrl)

	request, err := NewGetRequest("https://golang.org/")
	if err != nil {
		t.Fatalf("fail to create request: %v", err)
	}

	const requestNum = 100

	inputPipeline := func() chan *Request {
		output := make(chan *Request)
		go func() {
			defer close(output)
			for i := 0; i < requestNum; i++ {
				output <- request
			}
		}()
		return output
	}

	requestStrategyMock.EXPECT().CheckRestriction().Return(true).Times(requestNum)

	requestStrategyMock.EXPECT().Resource(
		gomock.AssignableToTypeOf(&Request{}),
	).DoAndReturn(
		func(r *Request) (string, error) { return r.URLHost(), nil },
	).Times(requestNum)

	requestSemaphoreMock.EXPECT().Acquire(
		ctx, gomock.AssignableToTypeOf(""),
	).Do(
		func(ctx context.Context, resource string) error {
			if resource != "golang.org" {
				t.Fatalf("expected golang.org, but got %s.\n", resource)
			}
			return nil
		},
	).Times(requestNum)

	requestSemaphoreMock.EXPECT().Release(
		gomock.AssignableToTypeOf(""),
	).Do(
		func(resource string) {
			// do nothing
		},
	).Times(requestNum)

	httpClientMock.EXPECT().Do(
		gomock.AssignableToTypeOf(&http.Request{}),
	).DoAndReturn(
		func(r *http.Request) (*http.Response, error) {
			return &http.Response{Request: r}, nil
		},
	).Times(requestNum)

	worker := newWorker(
		nil,
		requestStrategyMock,
		requestSemaphoreMock,
		httpClientMock,
		StdoutLogger{},
		10,
		[]func(request *Request){},
		[]func(response *Response){},
	)

	returnValueChan, err := worker.doRequest(inputPipeline())

	if err != nil {
		t.Fatalf("fail to Worker#doRequest: %v", err)
	}

	for returnValue := range returnValueChan {
		if returnValue.Request.URL != "https://golang.org/" {
			t.Fatalf("expect request url %s, but got %s", "https://golang.org/", returnValue.Request.URL)
		}
	}
}

// TestWorker_doRequestSuccess2
// CheckRestriction: false
func TestWorker_doRequestNotCheckRestriction(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	requestStrategyMock := NewMockRequestRestrictionStrategy(ctrl)
	requestSemaphoreMock := NewMockRequestSemaphore(ctrl)
	httpClientMock := NewMockHTTPClient(ctrl)

	request, err := NewGetRequest("https://golang.org/")
	if err != nil {
		t.Fatalf("fail to create request: %v", err)
	}

	const requestNum = 100

	inputPipeline := func() chan *Request {
		output := make(chan *Request)
		go func() {
			defer close(output)
			for i := 0; i < requestNum; i++ {
				output <- request
			}
		}()
		return output
	}

	requestStrategyMock.EXPECT().CheckRestriction().Return(false).Times(requestNum)

	httpClientMock.EXPECT().Do(
		gomock.AssignableToTypeOf(&http.Request{}),
	).DoAndReturn(
		func(r *http.Request) (*http.Response, error) {
			return &http.Response{Request: r}, nil
		},
	).Times(requestNum)

	worker := newWorker(
		nil,
		requestStrategyMock,
		requestSemaphoreMock,
		httpClientMock,
		StdoutLogger{},
		10,
		[]func(request *Request){},
		[]func(response *Response){},
	)

	returnValueChan, err := worker.doRequest(inputPipeline())

	if err != nil {
		t.Fatalf("fail to Worker#doRequest: %v", err)
	}

	for returnValue := range returnValueChan {
		if returnValue.Request.URL != "https://golang.org/" {
			t.Fatalf("expect request url %s, but got %s", "https://golang.org/", returnValue.Request.URL)
		}
	}
}

func TestWorker_doRequestRetry(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()

	requestStrategyMock := NewMockRequestRestrictionStrategy(ctrl)
	requestSemaphoreMock := NewMockRequestSemaphore(ctrl)
	httpClientMock := NewMockHTTPClient(ctrl)
	workerQueueMock := NewMockWorkerQueue(ctrl)
	loggerMock := NewMockLogger(ctrl)

	const requestURL = "https://golang.org/"

	request, err := NewGetRequest(requestURL)
	if err != nil {
		t.Fatalf("fail to create request: %v", err)
	}

	const requestNum = 100

	inputPipeline := func() chan *Request {
		output := make(chan *Request)
		go func() {
			defer close(output)
			for i := 0; i < requestNum; i++ {
				output <- request
			}
		}()
		return output
	}

	requestStrategyMock.EXPECT().CheckRestriction().Return(true).Times(requestNum)

	requestStrategyMock.EXPECT().ChangePriorityWhenRestricted(
		gomock.AssignableToTypeOf(&Request{}),
	).Do(
		func(r *Request) {
			// do nothing
		},
	).Times(requestNum)

	requestStrategyMock.EXPECT().Resource(
		gomock.AssignableToTypeOf(&Request{}),
	).DoAndReturn(
		func(r *Request) (string, error) { return r.URLHost(), nil },
	).Times(requestNum)

	requestSemaphoreMock.EXPECT().Acquire(
		ctx, gomock.AssignableToTypeOf(""),
	).DoAndReturn(
		func(ctx context.Context, resource string) error {
			if resource != "golang.org" {
				t.Fatalf("expected golang.org, but got %s.\n", resource)
			}
			// for retry
			return errors.New("retry request")
		},
	).Times(requestNum)

	workerQueueMock.EXPECT().RetryRequest(
		gomock.AssignableToTypeOf(&Request{}),
	).DoAndReturn(
		func(r *Request) error {
			return errors.New("")
		},
	).Times(requestNum)

	loggerMock.EXPECT().Infof("retry %s because worker failed to acquire resource.", requestURL).Times(requestNum)
	loggerMock.EXPECT().Errorf("fail to retry %s. this request is lost.").Times(requestNum)

	worker := newWorker(
		workerQueueMock,
		requestStrategyMock,
		requestSemaphoreMock,
		httpClientMock,
		loggerMock,
		10,
		[]func(request *Request){},
		[]func(response *Response){},
	)

	returnValueChan, err := worker.doRequest(inputPipeline())

	if err != nil {
		t.Fatalf("fail to Worker#doRequest: %v", err)
	}

	for returnValue := range returnValueChan {
		t.Fatalf("expect no response, got %s", returnValue.Request.URL)
	}
}
