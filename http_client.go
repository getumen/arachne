package lucy

//go:generate mockgen -source=$GOFILE -destination=mock_$GOFILE -package=$GOPACKAGE -self_package=github.com/getumen/lucy

import "net/http"

// HTTPClient is like http.Client
type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}
