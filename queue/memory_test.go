package queue

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/getumen/arachne"
	"github.com/wangjia184/sortedset"
)

func TestMemoryWorkerQueue_SubscribeRequests(t *testing.T) {
	ctx := context.Background()
	ctx, cancelFunc := context.WithCancel(ctx)

	q := memoryWorkerQueue{
		queue: sortedset.New(),
	}
	num := 1000

	go func() {
		for i := 0; i < num; i++ {
			r, _ := arachne.NewGetRequest(fmt.Sprintf("https://golang.org/%d", i))
			cond.L.Lock()
			q.queue.AddOrUpdate(r.URL, sortedset.SCORE(r.Priority), r)
			cond.Signal()
			cond.L.Unlock()
		}
		cancelFunc()
	}()

	counter := 0

	ch, err := q.SubscribeRequests(ctx)
	if err != nil {
		t.Fatalf("fail to subscribe")
	}
	for request := range ch {
		if request.URLHost() != "golang.org" {
			t.Fatalf("url domain mismatch")
		}
		counter++
	}
	cond.L.Lock()
	if counter+q.queue.GetCount() != num {
		t.Fatalf("expected %d, but got %d. some request missing.", num, counter+q.queue.GetCount())
	}
	cond.L.Unlock()
}

func TestMemoryWorkerQueue_SubscribeRequestsSlowPublication(t *testing.T) {
	ctx := context.Background()
	ctx, cancelFunc := context.WithCancel(ctx)

	q := memoryWorkerQueue{
		queue: sortedset.New(),
	}
	num := 1000

	go func() {
		for i := 0; i < num; i++ {
			r, _ := arachne.NewGetRequest(fmt.Sprintf("https://golang.org/%d", i))
			cond.L.Lock()
			q.queue.AddOrUpdate(r.URL, sortedset.SCORE(r.Priority), r)
			cond.Signal()
			cond.L.Unlock()
			time.Sleep(1 * time.Microsecond)
		}
		cancelFunc()
	}()

	counter := 0

	ch, err := q.SubscribeRequests(ctx)
	if err != nil {
		t.Fatalf("fail to subscribe")
	}
	for request := range ch {
		if request.URLHost() != "golang.org" {
			t.Fatalf("url domain mismatch")
		}
		counter++
	}
	cond.L.Lock()
	if counter+q.queue.GetCount() != num {
		t.Fatalf("expected %d, but got %d. some request missing.", num, counter+q.queue.GetCount())
	}
	cond.L.Unlock()
}

func TestMemoryWorkerQueue_SubscribeRequestsNoPublication(t *testing.T) {
	ctx := context.Background()
	ctx, cancelFunc := context.WithCancel(ctx)

	q := memoryWorkerQueue{
		queue: sortedset.New(),
	}
	num := 0

	go func() {
		for i := 0; i < num; i++ {
			r, _ := arachne.NewGetRequest(fmt.Sprintf("https://golang.org/%d", i))
			cond.L.Lock()
			q.queue.AddOrUpdate(r.URL, sortedset.SCORE(r.Priority), r)
			cond.Signal()
			cond.L.Unlock()
		}
		cancelFunc()
	}()

	counter := 0

	ch, err := q.SubscribeRequests(ctx)
	if err != nil {
		t.Fatalf("fail to subscribe")
	}
	for request := range ch {
		if request.URLHost() != "golang.org" {
			t.Fatalf("url domain mismatch")
		}
		counter++
	}
	cond.L.Lock()
	if counter+q.queue.GetCount() != num {
		t.Fatalf("expected %d, but got %d. some request missing.", num, counter+q.queue.GetCount())
	}
	cond.L.Unlock()
}

func TestMemoryWorkerQueue_RetryRequest(t *testing.T) {

	q := memoryWorkerQueue{
		queue: sortedset.New(),
	}
	wg := sync.WaitGroup{}
	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 100; j++ {
				r, _ := arachne.NewGetRequest(fmt.Sprintf("https://golang.org/%d", j))
				q.RetryRequest(r)
			}
		}()
	}

	wg.Wait()
	if q.queue.GetCount() != 100 {
		t.Fatalf("expected %d, but got %d", 100, q.queue.GetCount())
	}
}

func TestMemoryWorkerQueue_PublishRequest(t *testing.T) {

	q := memoryWorkerQueue{
		queue: sortedset.New(),
	}
	wg := sync.WaitGroup{}
	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 100; j++ {
				r, _ := arachne.NewGetRequest(fmt.Sprintf("https://golang.org/%d", j))
				q.PublishRequest(r)
			}
		}()
	}

	wg.Wait()
	if q.queue.GetCount() != 100 {
		t.Fatalf("expected %d, but got %d", 100, q.queue.GetCount())
	}
}
