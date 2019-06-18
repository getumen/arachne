//go:generate mockgen -source=$GOFILE -destination=mock_$GOFILE -package=$GOPACKAGE -self_package=github.com/getumen/lucy
package lucy

import "context"

type WorkerQueue interface {
	SubscribeRequests(ctx context.Context) (<-chan Request, error)
}
