package jsonapi

import (
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewClient(t *testing.T) {
	client, err := NewClient("http://example.org")
	assert.NoError(t, err)
	assert.NotNil(t, client)
	assert.Equal(t, "example.org", client.base.Host)

	client, err = NewClient("\n")
	assert.Error(t, err)
	assert.Nil(t, client)
}

type MockDoer struct {
	t          *testing.T
	statusCode int
	json       string
	req        *http.Request
	reqBody    string
}

func (m *MockDoer) Do(req *http.Request) (*http.Response, error) {
	m.req = req

	defer req.Body.Close()
	b, err := io.ReadAll(req.Body)
	if err != nil {
		return nil, err
	}
	m.reqBody = string(b)

	// m.t.Logf("got request: %+v", req)
	resp := &http.Response{
		StatusCode: m.statusCode,
		Body:       io.NopCloser(strings.NewReader(m.json)),
	}

	return resp, nil
}

func TestCall(t *testing.T) {
	m := &MockDoer{t: t, statusCode: 200, json: `{"baz":"quux"}`}

	client, err := NewClientWithDoer("http://example.org", m)
	assert.NoError(t, err)
	assert.NotNil(t, client)

	body := &struct {
		Foo string
	}{
		Foo: "bar",
	}
	result := &struct {
		Baz string
	}{}

	err = client.Call("/api/foo", body, result)
	assert.NoError(t, err)
	assert.NotNil(t, m.req)
	assert.Equal(t, "http://example.org/api/foo", m.req.URL.String())
	assert.Equal(t, http.MethodPost, m.req.Method)
	assert.Equal(t, "quux", result.Baz)

	// m = &MockDoer{t: t, statusCode: 200, json: `{"baz":"quux"}`}

	err = client.Call("\n", nil, nil)
	assert.Error(t, err)

	err = client.Call("/api/foo", func() {}, result)
	assert.Error(t, err)

	err = client.Call("/api/foo", nil, func() {})
	assert.Error(t, err)

	m = &MockDoer{t: t, statusCode: 400, json: `{}`}

	client, err = NewClientWithDoer("http://example.org", m)
	assert.NoError(t, err)
	assert.NotNil(t, client)
	err = client.Call("/api/foo", nil, result)
	assert.IsType(t, &HTTPStatusError{}, err)

}
