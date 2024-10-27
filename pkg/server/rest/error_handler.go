package rest

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/bioform/go-web-app-template/pkg/server/rest/apierror"
)

// HandleFunc type for functions with the same signature as Hello.Handle
type HandleFunc func(w http.ResponseWriter, r *http.Request) error

func ServeError(h HandleFunc) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				handleError(w, r, http.StatusInternalServerError, fmt.Errorf("panic occurred: %v", err))
			}
		}()

		if err := h(w, r); err != nil {
			handleError(w, r, http.StatusInternalServerError, err)
		}
	})
}

func handleError(w http.ResponseWriter, r *http.Request, statusCode int, err error) {
	apiErr, ok := err.(apierror.Error)
	if !ok {
		apiErr = apierror.NewWithError(statusCode, "Internal Server Error", err).(apierror.Error)
	}

	Encode(w, r, statusCode, apiErr)

	// if the error cause exists, log it
	if apiErr.Unwrap() != nil {
		slog.Error("error handling request", "error", apiErr)
	}
}
