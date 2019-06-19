package lucy

//go:generate mockgen -source=$GOFILE -destination=mock_$GOFILE -package=$GOPACKAGE -self_package=github.com/getumen/lucy

import (
	"net/url"

	"golang.org/x/xerrors"
)

// RequestRestrictionStrategy is an interface for request restriction. For example, request rate limit per domain.
type RequestRestrictionStrategy interface {
	CheckRestriction() bool
	Resource(request *Request) (string, error)
	ChangePriorityWhenRestricted(request *Request)
}

// DomainRestrictionStrategy is for restricting request num per domain.
type DomainRestrictionStrategy struct {
	DecreasePriorityNum int64
}

// CheckRestriction is always true.
func (*DomainRestrictionStrategy) CheckRestriction() bool {
	return true
}

// Resource returns domain name of request url.
func (*DomainRestrictionStrategy) Resource(request *Request) (string, error) {
	urlValue, err := url.ParseRequestURI(request.URL)
	if err != nil {
		return "", xerrors.Errorf("request url is invalid, but this never happens: %w", err)
	}
	return urlValue.Host, nil
}

// ChangePriorityWhenRestricted descrease the priority of request by DecreasePriorityNum
func (d *DomainRestrictionStrategy) ChangePriorityWhenRestricted(request *Request) {
	// decrease priority
	request.Priority += d.DecreasePriorityNum
}
