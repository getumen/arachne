package lucy

import (
	"net/url"
	"time"
)
import "golang.org/x/xerrors"

type Request struct {
	URL         string
	Method      string
	Headers     map[string]string
	Body        []byte
	Cookie      map[string]string
	Encoding    string
	nextRequest time.Time
	lastRequest time.Time
	stats       map[string]float64
	namespace   string
	meta        map[string]interface{}
}

type Response struct {
	Status  int
	Headers map[string]string
	Body    []byte
	Request *Request
}

func (r *Response) Follow(urlString string) (string, error) {
	requestUrl, err := url.Parse(r.Request.URL)
	if err != nil {
		return "", xerrors.Errorf("request url is invalid. this will never happen.: %w", err)
	}
	rawUrl, err := url.ParseRequestURI(urlString)
	if err != nil {
		return "", xerrors.Errorf("link url %s is invalid.: %w", urlString, err)
	}
	if rawUrl.Host == "" {
		rawUrl.Host = requestUrl.Host
		rawUrl.Scheme = requestUrl.Scheme
	}
	return rawUrl.String(), nil
}
