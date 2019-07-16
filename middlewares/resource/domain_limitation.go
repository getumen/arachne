package resource

import (
	"log"
	"sync"

	"github.com/getumen/arachne"
)

var (
	mutex         sync.Mutex
	domainCounter map[string]int64
)

func init() {
	mutex = sync.Mutex{}
	domainCounter = map[string]int64{}
}

// InMemoryDomainCounter is an in-memory domain counter
type InMemoryDomainCounter struct {
	maxRequestCount int64
}

// NewInMemoryDomainCounter is the InMemoryDomainCounter constructor
func NewInMemoryDomainCounter(maxRequestCount int64) *InMemoryDomainCounter {
	return &InMemoryDomainCounter{
		maxRequestCount: maxRequestCount,
	}
}

// RequestMiddleware is request middleware
func (c *InMemoryDomainCounter) RequestMiddleware(request *arachne.Request) {
	mutex.Lock()
	defer mutex.Unlock()
	count, ok := domainCounter[request.URLHost()]
	if ok {
		if count >= c.maxRequestCount {
			request.Meta["retry"] = true
			return
		}
	}
	count++
	domainCounter[request.URLHost()] = count
}

// ResponseMiddleware is response middleware
func (c *InMemoryDomainCounter) ResponseMiddleware(response *arachne.Response) {
	if retry, ok := response.Request.Meta["retry"]; ok {
		if retryFlag, ok := retry.(bool); retryFlag && ok {
			return
		}
	}
	mutex.Lock()
	defer mutex.Unlock()
	count, ok := domainCounter[response.Request.URLHost()]
	if ok {
		if count > 0 {
			count--
			domainCounter[response.Request.URLHost()] = count
		} else {
			// this never happened
			log.Panicf("the InMemoryDomainCounter is broken. count cannot be negative.")
		}
	}
}
