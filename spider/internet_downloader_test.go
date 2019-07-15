package spider

import (
	"testing"

	"github.com/getumen/lucy"
)

func TestDownloaDInternet_SimpleResponse(t *testing.T) {
	response := &lucy.Response{}
	response.Headers = map[string][]string{"Content-Type": {"text/html;utf-8"}}
	response.Body = []byte("<html><head><title>Test</title></head><body><a href='/doc/'>Documentation</a></body></html>")
	response.Request = &lucy.Request{}
	response.Request.URL = "https://golang.org/"

	requestList, err := DownloadInternet(response)
	if err != nil {
		t.Fatalf("fail to DownloadInternet: %v", err)
	}
	if len(requestList) != 1 {
		t.Fatalf("fail to parse html")
	}
	expected := "https://golang.org/doc/"
	if requestList[0].URL != expected {
		t.Fatalf("expected %s, but got %s", expected, requestList[0].URL)
	}
}
