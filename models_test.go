package lucy

import "testing"

func TestResponse_Follow(t *testing.T) {
	validResponse := Response{
		200,
		map[string]string{},
		[]byte{},
		&Request{
			"https://golang.org/",
			"GET",
			map[string]string{},
			[]byte{},
			map[string]string{},
			"",
			0,
			"default",
			map[string]interface{}{},
		}}

	invalidResponse := Response{
		200,
		map[string]string{},
		[]byte{},
		&Request{
			"gopher",
			"GET",
			map[string]string{},
			[]byte{},
			map[string]string{},
			"",
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
		if actualString, actualError := tt.response.Follow(tt.link);
			actualString != tt.expectedString || ((actualError != nil) != tt.expectedIsError) {
			t.Fatalf("test case %d: expectedString = %s, got = %s and error %v",
				i, tt.expectedString, actualString, actualError)
		}
	}
}
