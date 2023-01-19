package sensorpush

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestNewClient(t *testing.T) {
	client, err := NewClient("", "")
	assert.NoError(t, err)
	assert.NotNil(t, client)
}

type MockApiClient struct {
	mock.Mock
	t      *testing.T
	header http.Header
}

func (m *MockApiClient) Header() http.Header {
	if m.header == nil {
		m.header = make(http.Header)
	}

	return m.header
}

func (m *MockApiClient) Call(url string, body interface{}, response interface{}) error {
	args := m.MethodCalled("Call", url, body, response)
	return args.Error(0)
}

// TODO: have auth/access failure depending on USER/PASS?  Have a retry/count
// test?
func (m *MockApiClient) SetUpAuth() *mock.Call {
	return m.
		On("Call",
			"oauth/authorize",
			mock.AnythingOfType("*sensorpush.authorizeRequest"),
			mock.AnythingOfType("*sensorpush.authorizeResponse")).
		Run(func(args mock.Arguments) {
			r := args.Get(2).(*authorizeResponse)
			r.Authorization = "AUTH-TOKEN"
		}).
		Return(nil).
		On("Call",
			"oauth/accesstoken",
			mock.AnythingOfType("*sensorpush.accessTokenRequest"),
			mock.AnythingOfType("*sensorpush.accessTokenResponse")).
		Run(func(args mock.Arguments) {
			r := args.Get(2).(*accessTokenResponse)
			r.AccessToken = "ACCESS-TOKEN"
		}).
		Return(nil)
}

// creates a sensorpush Client with a mocked API handler
func NewMockedClient(t *testing.T) (*Client, *mock.Call) {
	m := &MockApiClient{t: t}
	m.Test(t)

	client, err := NewClientWithApiClient("USER", "PASS", m)
	assert.NoError(t, err)
	assert.NotNil(t, client)

	return client, m.SetUpAuth()
}
