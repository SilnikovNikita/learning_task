package service

import "testing"

func TestService_Masking(t *testing.T) {
	tests := []struct {
		name     string
		input    []byte
		expected []byte
	}{
		// TODO: Add test cases.

		{
			name:     "simple URL replacement",
			input:    []byte("Visit https://example.com for more info"),
			expected: []byte("Visit https://*********** for more info"),
		},
		{
			name:     "URL at start",
			input:    []byte("https://site.org is great"),
			expected: []byte("https://******** is great"),
		},
		{
			name:     "URL at end",
			input:    []byte("Check out https://domain.net"),
			expected: []byte("Check out https://**********"),
		},
		{
			name:     "multiple URLs - multiple URLs",
			input:    []byte("Links: https://first.com and https://second.org"),
			expected: []byte("Links: https://********* and https://**********"),
		},
		{
			name:     "no URL in text",
			input:    []byte("Just regular text without URLs"),
			expected: []byte("Just regular text without URLs"),
		},
		{
			name:     "partial URL match",
			input:    []byte("This http://not.https.com won't match"),
			expected: []byte("This http://not.https.com won't match"),
		},
		{
			name:     "empty input",
			input:    []byte(""),
			expected: []byte(""),
		},
		{
			name:     "URL with path",
			input:    []byte("Go to https://site.com/path/page.html now"),
			expected: []byte("Go to https://*********************** now"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Service{}
			inputCopy := make([]byte, len(tt.input))
			copy(inputCopy, tt.input)
			s.Masking(inputCopy)

			if string(inputCopy) != string(tt.expected) {
				t.Errorf("UrlToAsterisk() = %q, want %q", inputCopy, tt.expected)
			}
		})
	}
}
