package lucy

import "context"

type WorkerQueue interface {
	SubscribeRequests(ctx context.Context) (<-chan Request, error)
}
