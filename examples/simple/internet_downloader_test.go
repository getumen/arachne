package main

import (
	"context"
	"log"
	"net/http"
	"testing"
	"time"

	"github.com/getumen/lucy"
	"github.com/getumen/lucy/logger"
	"github.com/getumen/lucy/middlewares/resource"
	"github.com/getumen/lucy/queue"
	"github.com/getumen/lucy/spider"
)

func TestSimpleCrawler(t *testing.T) {
	builder := lucy.NewWorkerBuilder()
	builder.SetLogger(logger.NewStdoutLogger(lucy.DebugLevel))
	builder.SetHTTPClient(&http.Client{})
	queue, err := queue.NewMemoryWorkerQueue()
	if err != nil {
		log.Fatalf("fail to create queue: %v", err)
	}
	builder.SetWorkerQueue(queue)
	builder.SetSpider(spider.DownloadInternet)

	worker, err := builder.Build()
	if err != nil {
		log.Fatalf("fail to create worker: %v", err)
	}

	domainRestriction := resource.NewInMemoryDomainCounter(1)

	builder.AddRequestMiddleware(domainRestriction.RequestMiddleware)
	builder.AddResponseMiddleware(domainRestriction.ResponseMiddleware)

	workerRestriction := resource.NewRequestCounter(1)
	builder.AddRequestMiddleware(workerRestriction.RequestMiddleware)
	builder.AddResponseMiddleware(workerRestriction.ResponseMiddleware)

	ctx := context.Background()

	ctx, cancelFunc := context.WithCancel(ctx)

	go func() {
		time.Sleep(3 * time.Second)
		cancelFunc()
	}()

	worker.StartWithFirstRequest(ctx, "http://example.com/")
}
