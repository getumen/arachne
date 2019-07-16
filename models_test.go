package arachne

import (
	"net/http"
	"testing"
)

func TestResponse_Follow(t *testing.T) {
	request, err := NewGetRequest("https://golang.org/")
	if err != nil {
		t.Fatalf("fail to create request")
	}
	validResponse := Response{
		200,
		http.Header{},
		[]byte{},
		request}

	request, err = NewGetRequest("gopher")
	if err != nil {
		t.Fatalf("fail to create request")
	}
	invalidResponse := Response{
		200,
		http.Header{},
		[]byte{},
		request}

	tests := []struct {
		response        Response
		link            string
		expectedString  string
		expectedIsError bool
	}{
		{validResponse, "/doc/", "https://golang.org/doc/", false},
		{validResponse, "$$", "", true},
		{invalidResponse, "/doc/", "", true},
	}

	for i, tt := range tests {
		if actualString, actualError := tt.response.Follow(tt.link); actualString != tt.expectedString || ((actualError != nil) != tt.expectedIsError) {
			t.Fatalf("test case %d: expectedString = %s, got = %s and error %v",
				i, tt.expectedString, actualString, actualError)
		}
	}
}
