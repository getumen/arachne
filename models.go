package lucy

import (
	"fmt"
	"net/url"
)
import "golang.org/x/xerrors"

type Request struct {
	URL       string
	Method    string
	Headers   map[string]string
	Body      []byte
	Cookie    map[string]string
	Encoding  string
	Priority  int64
	QueueName string
	Meta      map[string]interface{}
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
		return "", xerrors.Errorf("request url %s is invalid. this will never happen.: %w", r.Request.URL, err)
	} else if requestUrl.Host == "" || requestUrl.Scheme == "" {
		return "", xerrors.New(fmt.Sprintf("request url %s is invalid. this will never happen.", r.Request.URL))
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
