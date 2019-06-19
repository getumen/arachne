package lucy

import (
	"fmt"
	"net/url"

	"golang.org/x/xerrors"
)

// Request is a domain model that represents http request.
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

// Response is a domain model that represents http response.
type Response struct {
	Status  int
	Headers map[string]string
	Body    []byte
	Request *Request
}

// NewGetRequest creates simple GET request.
func NewGetRequest(urlStr string) *Request {
	return &Request{
		URL:       urlStr,
		Method:    "GET",
		Headers:   map[string]string{},
		Body:      []byte{},
		Cookie:    map[string]string{},
		Encoding:  "utf-8",
		Priority:  0,
		QueueName: "default",
		Meta:      map[string]interface{}{},
	}
}

// Follow creates a url whose url schema and host is the same as those of response.
func (r *Response) Follow(urlString string) (string, error) {
	requestURL, err := url.Parse(r.Request.URL)
	if err != nil {
		return "", xerrors.Errorf("request url %s is invalid. this will never happen.: %w", r.Request.URL, err)
	} else if requestURL.Host == "" || requestURL.Scheme == "" {
		return "", xerrors.New(fmt.Sprintf("request url %s is invalid. this will never happen.", r.Request.URL))
	}
	rawURL, err := url.ParseRequestURI(urlString)
	if err != nil {
		return "", xerrors.Errorf("link url %s is invalid.: %w", urlString, err)
	}
	if rawURL.Host == "" {
		rawURL.Host = requestURL.Host
		rawURL.Scheme = requestURL.Scheme
	}
	return rawURL.String(), nil
}

// FollowRequest creates a simple GET request whose the schema and the host of the url is the same as those of response.
func (r *Response) FollowRequest(urlString string) (*Request, error) {
	requestURL, err := r.Follow(urlString)
	if err != nil {
		return nil, xerrors.Errorf("fail to make request url.: %w", err)
	}
	return NewGetRequest(requestURL), nil
}
