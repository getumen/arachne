package main

import (
	"context"
	"log"
	"net/http"
	"testing"
	"time"

	"github.com/getumen/arachne"
	"github.com/getumen/arachne/builder"
	"github.com/getumen/arachne/logger"
	"github.com/getumen/arachne/middlewares/resource"
	"github.com/getumen/arachne/queue"
	"github.com/getumen/arachne/spider"
)

func TestSimpleCrawler(t *testing.T) {
	workerBuilder := builder.NewWorkerBuilder()
	workerBuilder.SetLogger(logger.NewStdoutLogger(arachne.InfoLevel))
	httpClient := &http.Client{}
	httpClient.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}
	workerBuilder.SetHTTPClient(httpClient)
	memoryQueue, err := queue.NewMemoryWorkerQueue()
	if err != nil {
		log.Fatalf("fail to create queue: %v", err)
	}
	workerBuilder.SetWorkerQueue(memoryQueue)
	workerBuilder.SetSpider(spider.DownloadInternet)

	worker, err := workerBuilder.Build()
	if err != nil {
		log.Fatalf("fail to create worker: %v", err)
	}

	domainRestriction := resource.NewInMemoryDomainCounter(1)

	worker.RequestMiddlewares = append(worker.RequestMiddlewares, domainRestriction.RequestMiddleware)
	worker.ResponseMiddlewares = append(worker.ResponseMiddlewares, domainRestriction.ResponseMiddleware)

	workerRestriction := resource.NewRequestCounter(1)
	worker.RequestMiddlewares = append(worker.RequestMiddlewares, workerRestriction.RequestMiddleware)
	worker.ResponseMiddlewares = append(worker.ResponseMiddlewares, workerRestriction.ResponseMiddleware)

	worker.RequestMiddlewares = append(worker.RequestMiddlewares, worker.RetryMiddleware)

	ctx := context.Background()

	ctx, cancelFunc := context.WithCancel(ctx)

	go func() {
		time.Sleep(3 * time.Second)
		cancelFunc()
	}()

	worker.StartWithFirstRequest(ctx, "http://example.com/")
}
