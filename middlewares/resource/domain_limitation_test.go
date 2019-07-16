package resource

import (
	"sync"
	"sync/atomic"
	"testing"

	"github.com/getumen/arachne"
)

func TestInMemoryDomainCounter_RequestMiddleware(t *testing.T) {
	var maxRequestCount int64 = 10
	var loop int64 = 1000
	var retryCount int64
	var i int64

	target := NewInMemoryDomainCounter(maxRequestCount)
	wg := sync.WaitGroup{}
	for i = 0; i < loop; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			request, _ := arachne.NewGetRequest("https://golang.org/")
			target.RequestMiddleware(request)
			if retry, ok := request.Meta["retry"]; ok {
				if retryFlag, ok := retry.(bool); ok && retryFlag {
					atomic.AddInt64(&retryCount, 1)
				}
			}
		}()
	}
	wg.Wait()
	if retryCount != loop-maxRequestCount {
		t.Fatalf("expected %d, but got %d", loop-maxRequestCount, retryCount)
	}
}

func TestInMemoryDomainCounter_ResponseMiddleware(t *testing.T) {
	var maxRequestCount int64 = 1000
	var loop int64 = 1000
	var i int64

	target := NewInMemoryDomainCounter(maxRequestCount)
	domainCounter["golang.org"] = 1000
	wg := sync.WaitGroup{}
	for i = 0; i < loop; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			request, _ := arachne.NewGetRequest("https://golang.org/")
			response := &arachne.Response{Request: request}
			target.ResponseMiddleware(response)
		}()
	}
	wg.Wait()
	if domainCounter["golang.org"] != 0 {
		t.Fatalf("expected %d, but got %d", 0, domainCounter["golang.org"])
	}
}
