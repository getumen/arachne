package lucy

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"github.com/golang/mock/gomock"
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
	loggerMock := NewMockLogger(ctrl)

	workerQueueMock.EXPECT().SubscribeRequests(ctx).Return(subscribeChan, nil)

	worker := newWorker(
		workerQueueMock,
		nil,
		loggerMock,
		[]func(request *Request){}, []func(response *Response){},
		nil,
	)
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
	loggerMock := NewMockLogger(ctrl)

	workerQueueMock.EXPECT().SubscribeRequests(ctx).Return(nil, xerrors.New("some error."))

	worker := newWorker(
		workerQueueMock,
		nil,
		loggerMock,
		[]func(request *Request){}, []func(response *Response){},
		nil,
	)
	reqChan, err := worker.subscribe(ctx)

	if reqChan != nil || err == nil {
		t.Fatalf("fail to propagate error.")
	}
}

// TestWorker_doRequestSuccess
func TestWorker_doRequestRestrictionByDomain(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	httpClientMock := NewMockHTTPClient(ctrl)
	loggerMock := NewMockLogger(ctrl)

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

	httpClientMock.EXPECT().Do(
		gomock.AssignableToTypeOf(&http.Request{}),
	).DoAndReturn(
		func(r *http.Request) (*http.Response, error) {
			return &http.Response{Request: r}, nil
		},
	).Times(requestNum)

	worker := newWorker(
		nil,
		httpClientMock,
		loggerMock,
		[]func(request *Request){},
		[]func(response *Response){},
		nil,
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

func TestWorker_applySpiderReturnNewRequests(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	loggerMock := NewMockLogger(ctrl)

	const requestURL = "https://golang.org/"

	request, _ := NewGetRequest(requestURL)

	worker := newWorker(
		nil,
		nil,
		loggerMock,
		[]func(request *Request){},
		[]func(response *Response){},
		func(response *Response) ([]*Request, error) {
			return []*Request{
				request,
				request,
				request,
			}, nil
		},
	)

	const responseNum = 100

	inputPipeline := func() chan *Response {
		output := make(chan *Response)
		go func() {
			defer close(output)
			for i := 0; i < responseNum; i++ {
				output <- &Response{}
			}
		}()
		return output
	}

	returnValueChan, err := worker.applySpider(inputPipeline())
	if err != nil {
		t.Fatalf("expect err == nil, but got %v", err)
	}

	requestCounter := 0
	for returnValue := range returnValueChan {
		if returnValue.URL != requestURL {
			t.Fatalf("expect no response, got %s", returnValue.URL)
		}
		requestCounter++
	}
	if requestCounter != 300 {
		t.Fatalf("expect requestCounter == 300, but got %d", requestCounter)
	}
}

func TestWorker_applySpiderReturnError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	loggerMock := NewMockLogger(ctrl)

	worker := newWorker(
		nil,
		nil,
		loggerMock,
		[]func(request *Request){},
		[]func(response *Response){},
		func(response *Response) ([]*Request, error) {
			return nil, errors.New("error")
		},
	)

	const responseNum = 100

	loggerMock.EXPECT().Infof(
		"spider error: %v", gomock.AssignableToTypeOf(errors.New("")),
	).Times(responseNum)

	inputPipeline := func() chan *Response {
		output := make(chan *Response)
		go func() {
			defer close(output)
			for i := 0; i < responseNum; i++ {
				output <- &Response{}
			}
		}()
		return output
	}

	returnValueChan, err := worker.applySpider(inputPipeline())

	if err != nil {
		t.Fatalf("expect err == nil, but got %v", err)
	}

	requestCounter := 0
	for returnValue := range returnValueChan {
		if returnValue != nil {
			t.Fatalf("expect no result, but got %v", returnValue)
		}
		requestCounter++
	}
	if requestCounter != 0 {
		t.Fatalf("expect requestCounter == 0, but got %d", requestCounter)
	}
}

func TestWorker_publishRequestPublish(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	loggerMock := NewMockLogger(ctrl)
	workerQueueMock := NewMockWorkerQueue(ctrl)

	worker := newWorker(
		workerQueueMock,
		nil,
		loggerMock,
		[]func(request *Request){},
		[]func(response *Response){},
		nil,
	)

	const requestNum = 100

	inputPipeline := func() chan *Request {
		output := make(chan *Request)
		go func() {
			defer close(output)
			for i := 0; i < requestNum; i++ {
				output <- &Request{URL: "https://golang.org/"}
			}
		}()
		return output
	}

	workerQueueMock.EXPECT().PublishRequest(
		gomock.AssignableToTypeOf(&Request{}),
	).DoAndReturn(
		func(r *Request) error {
			if r.URL != "https://golang.org/" {
				t.Fatalf("expect https://golang.org/, but got %s", r.URL)
			}
			return nil
		},
	).Times(requestNum)

	_ = worker.publishRequest(inputPipeline())
}

func TestWorker_publishRequestFailToPublish(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	loggerMock := NewMockLogger(ctrl)
	workerQueueMock := NewMockWorkerQueue(ctrl)

	worker := newWorker(
		workerQueueMock,
		nil,
		loggerMock,
		[]func(request *Request){},
		[]func(response *Response){},
		nil,
	)

	const requestNum = 100

	inputPipeline := func() chan *Request {
		output := make(chan *Request)
		go func() {
			defer close(output)
			for i := 0; i < requestNum; i++ {
				output <- &Request{URL: "https://golang.org/"}
			}
		}()
		return output
	}

	workerQueueMock.EXPECT().PublishRequest(
		gomock.AssignableToTypeOf(&Request{}),
	).DoAndReturn(
		func(r *Request) error {
			if r.URL != "https://golang.org/" {
				t.Fatalf("expect https://golang.org/, but got %s", r.URL)
			}
			return errors.New("")
		},
	).Times(requestNum)

	loggerMock.EXPECT().Errorf("fail to publish request: %s", "https://golang.org/").Times(requestNum)

	_ = worker.publishRequest(inputPipeline())
}
