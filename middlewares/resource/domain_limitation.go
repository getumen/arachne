package resource

import (
	"sync"

	"github.com/getumen/lucy"
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
	maxRequestCounter int64
}

// RequestMiddleware is request middleware
func (c *InMemoryDomainCounter) RequestMiddleware(request *lucy.Request) {
	mutex.Lock()
	defer mutex.Unlock()
	count, ok := domainCounter[request.URLHost()]
	if ok {
		if count > c.maxRequestCounter {
			request.Meta["retry"] = true
			return
		}
	}
	count++
	domainCounter[request.URLHost()] = count
}

// ResponseMiddleware is response middleware
func (c *InMemoryDomainCounter) ResponseMiddleware(response *lucy.Response) {
	mutex.Lock()
	defer mutex.Unlock()
	count, ok := domainCounter[response.Request.URLHost()]
	if ok {
		count--
		domainCounter[response.Request.URLHost()] = count
	}
}
