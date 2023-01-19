package jsonapi

import (
	"testing"
)

func TestError(t *testing.T) {
	err := HTTPStatusError{StatusCode: 200, Body: []byte("some text")}
	actual := err.Error()
	expected := "200 - OK"
	if actual != expected {
		t.Errorf("expected %q, got %q", expected, actual)
	}
}
