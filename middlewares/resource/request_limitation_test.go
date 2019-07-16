package resource

import (
	"testing"

	"github.com/getumen/arachne"
)

func TestRequestCounter_RequestMiddleware(t *testing.T) {
	var loop int64 = 1000
	target := NewRequestCounter(loop)

	for i := loop; i > 0; i-- {
		r, _ := arachne.NewGetRequest("https://golang.org/")
		target.RequestMiddleware(r)
	}
	// this returns error because of out of resource
	ok := target.sema.TryAcquire(1)
	if ok {
		t.Fatalf("expected resource does not exist, but got exists")
	}
}

func TestRequestCounter_ResponseMiddleware(t *testing.T) {
	var loop int64 = 1000
	target := NewRequestCounter(loop)
	target.sema.TryAcquire(loop)

	for i := loop; i > 0; i-- {
		r, _ := arachne.NewGetRequest("https://golang.org/")
		response := &arachne.Response{
			Request: r,
		}
		target.ResponseMiddleware(response)
	}
	// test counting semaphore is exactly zero
	defer func() {
		err := recover()
		if err != "semaphore: released more than held" {
			t.Errorf("got %v\nwant %v", err, "semaphore: released more than held")
		}
	}()
	// this returns error because of out of resource
	target.sema.Release(1)
}
