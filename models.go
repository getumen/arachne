package lucy

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"

	"golang.org/x/xerrors"
)

// Request is a domain model that represents http request.
type Request struct {
	URL        string
	Method     string
	Header     http.Header
	Body       []byte
	Priority   int64
	QueueName  string
	Meta       map[string]interface{}
	requestURL *url.URL
}

// NewGetRequest creates simple GET request.
func NewGetRequest(urlStr string) (*Request, error) {
	requestURL, err := url.Parse(urlStr)
	if err != nil {
		return nil, xerrors.Errorf("fail to make request url.: %w", err)
	}
	request := new(Request)
	request.URL = urlStr
	request.Method = "GET"
	request.Header = http.Header{}
	request.Body = []byte{}
	request.Priority = 0
	request.QueueName = "default"
	request.Meta = map[string]interface{}{}
	request.requestURL = requestURL
	return request, nil
}

// NewRequestFromHTTPRequest constructs Request from http.Request
func NewRequestFromHTTPRequest(request *http.Request) (*Request, error) {
	r := new(Request)
	r.URL = request.URL.String()
	r.requestURL = request.URL
	r.Method = request.Method
	if r.Header == nil {
		r.Header = map[string][]string{}
	}
	for key, values := range request.Header {
		for _, value := range values {
			r.Header.Add(key, value)
		}
	}
	if request.Body != nil {
		body, err := ioutil.ReadAll(request.Body)
		if err != nil {
			return nil, xerrors.Errorf("fail to read request body url: %s: %w ", request.URL.String(), err)
		}
		defer request.Body.Close()
		r.Body = body
	}
	r.QueueName = "default"
	r.Meta = map[string]interface{}{}
	return r, nil
}

// HTTPRequest constructs http.Request from Request
func (r *Request) HTTPRequest() (*http.Request, error) {
	o, err := http.NewRequest(r.Method, r.URL, r.BodyReader())
	if err != nil {
		return nil, xerrors.Errorf(
			"fail to create new request. URL(%s), Mehotd(%s) or Body(%s) are invalid.: %w ",
			r.URL, r.Method, string(r.Body), err)
	}
	for key, values := range r.Header {
		for _, value := range values {
			o.Header.Add(key, value)
		}
	}
	return o, nil
}

// URLHost returns the host of the request url.
func (r *Request) URLHost() string {
	return r.requestURL.Host
}

// BodyReader returns io.Reader of Body
func (r *Request) BodyReader() io.Reader {
	return bytes.NewBuffer(r.Body)
}

// Response is a domain model that represents http response.
type Response struct {
	StatusCode int
	Headers    http.Header
	Body       []byte
	Request    *Request
}

// NewResponseFromHTTPResponse constructs Response from http.Response
func NewResponseFromHTTPResponse(response *http.Response) (*Response, error) {
	r := new(Response)
	r.StatusCode = response.StatusCode
	r.Headers = response.Header
	if response.Body != nil {
		bodyBytes, err := ioutil.ReadAll(response.Body)
		if err != nil {
			return nil, xerrors.Errorf("fail to read body of  %s: %w ", response.Request.URL.String(), err)
		}
		defer response.Body.Close()
		r.Body = bodyBytes
	}
	request, err := NewRequestFromHTTPRequest(response.Request)
	if err != nil {
		return nil, xerrors.Errorf("fail to read body of  %s: %w ", response.Request.URL.String(), err)
	}
	r.Request = request
	return r, nil
}

// Follow creates a url whose url schema and host is the same as those of response.
func (r *Response) Follow(urlString string) (string, error) {
	requestURL, err := url.Parse(r.Request.URL)
	if err != nil {
		return "", xerrors.Errorf("request url %s is invalid. this will never happened: %w", r.Request.URL, err)
	} else if requestURL.Host == "" || requestURL.Scheme == "" {
		return "", xerrors.New(fmt.Sprintf("request url %s is invalid. this will be never happened", r.Request.URL))
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
	req, err := NewGetRequest(requestURL)
	if err != nil {
		return nil, xerrors.Errorf("fail to make request.: %w", err)
	}
	return req, nil
}

// ContentType returns detected Content-Type of the response body.
// If it cannot determine a more specific one, it returns "application/octet-stream".
func (r *Response) ContentType() string {
	return r.Headers.Get("Content-Type")
}

// Text returns string(Response.Body)
// Note that this method does not decode body.
// To ensure that the Response.Body is decoded,
// use http.Client with DefaultTransport.
func (r *Response) Text() string {
	return string(r.Body)
}
