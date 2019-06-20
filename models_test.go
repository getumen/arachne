package lucy

import (
	"net/http"
	"testing"
)

func TestResponse_Follow(t *testing.T) {
	validResponse := Response{
		200,
		http.Header{},
		[]byte{},
		&Request{
			"https://golang.org/",
			"GET",
			http.Header{},
			[]byte{},
			0,
			"default",
			map[string]interface{}{},
		}}

	invalidResponse := Response{
		200,
		http.Header{},
		[]byte{},
		&Request{
			"gopher",
			"GET",
			http.Header{},
			[]byte{},
			0,
			"default",
			map[string]interface{}{},
		}}

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
