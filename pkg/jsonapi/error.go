package jsonapi

import (
	"fmt"
	"net/http"
)

// HTTPStatusError represents an error from the HTTP layer (4xx, 5xx)
type HTTPStatusError struct {
	StatusCode int
	Body       []byte
}

func (e *HTTPStatusError) Error() string {
	return fmt.Sprintf("%d - %s", e.StatusCode, http.StatusText(e.StatusCode))
}
