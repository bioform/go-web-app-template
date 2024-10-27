package server

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRoutes(t *testing.T) {
	mux := RegisterRoutes()

	tests := []struct {
		method   string
		path     string
		expected int
	}{
		{"GET", "/", http.StatusOK},
		{"GET", "/health", http.StatusOK},
		{"POST", "/", http.StatusMethodNotAllowed},
		{"GET", "/nonexistent", http.StatusNotFound},
		{"POST", "/health", http.StatusMethodNotAllowed},
		{"PUT", "/", http.StatusMethodNotAllowed},
		{"DELETE", "/", http.StatusMethodNotAllowed},
		{"GET", "/anothernonexistent", http.StatusNotFound},
	}

	for _, tt := range tests {
		t.Run(tt.method+" "+tt.path, func(t *testing.T) {
			req, err := http.NewRequest(tt.method, tt.path, nil)
			if err != nil {
				t.Fatalf("could not create request: %v", err)
			}
			rr := httptest.NewRecorder()

			mux.ServeHTTP(rr, req)

			if rr.Code != tt.expected {
				t.Errorf("expected status %v; got %v for %v %v", tt.expected, rr.Code, tt.method, tt.path)
			}
		})
	}
}
