package apierror

import "net/http"

func InvalidRequestData(errors map[string]string) error {
	return New(http.StatusUnprocessableEntity, errors)
}

func RecordNotFound() error {
	return New(http.StatusNotFound, "Record not found")
}
