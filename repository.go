//go:generate mockgen -source=$GOFILE -destination=mock_$GOFILE -package=$GOPACKAGE
package lucy

import "context"

type WorkerQueue interface {
	SubscribeRequests(ctx context.Context) (<-chan Request, error)
}
