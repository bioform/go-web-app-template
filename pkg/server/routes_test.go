package server

import (
	"net/http"
	"net/http/httptest"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Route Tests", func() {
	var (
		router http.Handler
	)

	BeforeEach(func() {
		router = RegisterRoutes()
	})

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
		tt := tt // capture range variable
		It(tt.method+" "+tt.path, func() {
			req, err := http.NewRequest(tt.method, tt.path, nil)
			Expect(err).NotTo(HaveOccurred())

			rr := httptest.NewRecorder()
			router.ServeHTTP(rr, req)

			Expect(rr.Code).To(Equal(tt.expected))
		})
	}
})
