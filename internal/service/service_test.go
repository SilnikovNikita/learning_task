package service

import (
	"errors"
	mocks "task_1/internal/service/mocks"
	"testing"
)

func TestService_Masking(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		// TODO: Add test cases.

		{
			name:     "simple URL replacement",
			input:    "Visit https://example.com for more info",
			expected: "Visit https://*********** for more info",
		},
		{
			name:     "URL at start",
			input:    "https://site.org is great",
			expected: "https://******** is great",
		},
		{
			name:     "URL at end",
			input:    "Check out https://domain.net",
			expected: "Check out https://**********",
		},
		{
			name:     "multiple URLs - multiple URLs",
			input:    "Links: https://first.com and https://second.org",
			expected: "Links: https://********* and https://**********",
		},
		{
			name:     "no URL in text",
			input:    "Just regular text without URLs",
			expected: "Just regular text without URLs",
		},
		{
			name:     "partial URL match",
			input:    "This http://not.https.com won't match",
			expected: "This http://not.https.com won't match",
		},
		{
			name:     "empty input",
			input:    "",
			expected: "",
		},
		{
			name:     "URL with path",
			input:    "Go to https://site.com/path/page.html now",
			expected: "Go to https://*********************** now",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Service{}

			output := s.masking(tt.input)

			if output != tt.expected {
				t.Errorf("UrlToAsterisk() = %q, want %q", output, tt.expected)
			}
		})
	}
}

func TestService_Run(t *testing.T) {
	tests := []struct {
		name      string
		args      []string
		want      []string
		prodError error
		presError error
		wantErr   bool
	}{
		// TODO: Add test cases.
		{
			name: "success",
			args: []string{
				"https://example.com",
				"https://example.com123",
			},
			want: []string{
				"https://***********",
				"https://**************",
			},
		},
		{
			name: "empty",
			args: []string{},
			want: []string{},
		},
		{
			name:      "error produce",
			prodError: errors.New("produce error"),
			wantErr:   true,
		},
		{
			name:      "error present",
			args:      []string{"https://example.com"},
			want:      []string{"https://***********"},
			presError: errors.New("present error"),
			wantErr:   true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			prod := mocks.NewMockProducer(t)
			pres := mocks.NewMockPresenter(t)

			if tt.name != "empty" || len(tt.args) > 0 {
				prod.
					On("Produce").
					Once().
					Return(tt.args, tt.prodError)
			}

			if tt.prodError == nil && len(tt.want) > 0 {
				pres.
					On("Present", tt.want).
					Once().
					Return(tt.presError)
			}

			if tt.name == "empty" {
				prod.
					On("Produce").
					Once().
					Return(tt.args, nil)
				pres.
					On("Present", tt.want).
					Once().
					Return(nil)
			}

			s := NewService(prod, pres)

			if err := s.Run(); (err != nil) != tt.wantErr {
				t.Errorf("Run() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
