//go:generate mockgen -source=$GOFILE -destination=mock_$GOFILE -package=$GOPACKAGE -self_package=github.com/getumen/lucy
package lucy

import (
	"golang.org/x/xerrors"
	"net/url"
)

type RequestRestrictionStrategy interface {
	CheckRestriction() bool
	Resource(request *Request) (string, error)
	ChangePriorityWhenRestricted(request *Request)
}

type DomainRestrictionStrategy struct {
}

func (*DomainRestrictionStrategy) CheckRestriction() bool {
	return true
}

func (*DomainRestrictionStrategy) Resource(request *Request) (string, error) {
	urlValue, err := url.ParseRequestURI(request.URL)
	if err != nil {
		return "", xerrors.Errorf("request url is invalid, but this never happens: %w", err)
	}
	return urlValue.Host, nil
}

func (*DomainRestrictionStrategy) ChangePriorityWhenRestricted(request *Request) {
	// decrease priority
	request.Priority += 10
}
