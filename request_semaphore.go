//go:generate mockgen -source=$GOFILE -destination=mock_$GOFILE -package=$GOPACKAGE -self_package=github.com/getumen/lucy
package lucy

import "context"

type RequestSemaphore interface {
	Acquire(ctx context.Context, resource string) error
	Release(resource string)
}
