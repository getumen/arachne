package lucy

//go:generate mockgen -source=$GOFILE -destination=mock_$GOFILE -package=$GOPACKAGE -self_package=github.com/getumen/lucy

import "context"

// RequestSemaphore prevents lucy crawler from sending too many requests to the resources.
type RequestSemaphore interface {
	Acquire(ctx context.Context, resource string) error
	Release(resource string)
}
