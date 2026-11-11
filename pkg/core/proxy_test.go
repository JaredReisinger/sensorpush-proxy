package core

import (
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/jaredreisinger/sensorpush-proxy/pkg/sensorpush"
	"github.com/stretchr/testify/assert"
)

var dummySample = sensorpush.Sample{
	Altitude:           1,
	BarometricPressure: 2,
	Dewpoint:           3,
	Humidity:           4,
	Observed:           time.Date(2000, 1, 1, 12, 0, 0, 0, time.UTC),
	Temperature:        6,
	Vpd:                7,
}

type mockSPClient struct {
	calls int
}

func (c *mockSPClient) LastSample(id string) (*sensorpush.Sample, error) {
	c.calls++

	if id == "DUMMY" {
		return &dummySample, nil
	}

	return nil, errors.New("invalid sensor ID")
}

func newTestServer(client *mockSPClient) *Server {
	if client == nil {
		client = &mockSPClient{}
	}
	s, _ := newServer(
		"127.0.0.1:0",
		client,
		map[string]string{"dummy": "DUMMY"},
		250*time.Millisecond,
	)
	return s
}

func TestSensorsHandler(t *testing.T) {
	// Set up some bogus data...
	s := newTestServer(nil)

	s.lastSamples["DUMMY"] = dummySample

	req := httptest.NewRequest(http.MethodGet, "/sensors", nil)
	w := httptest.NewRecorder()
	s.sensorsHandler(w, req)
	res := w.Result()

	assert.Equal(t, http.StatusOK, res.StatusCode)
	assert.Equal(t, "application/json", res.Header.Get("Content-Type"))

	defer res.Body.Close()
	b, err := io.ReadAll(res.Body)
	assert.NoError(t, err)
	assert.Equal(t, `{"DUMMY":{"altitude":1,"barometric_pressure":2,"dewpoint":3,"humidity":4,"observed":"2000-01-01T12:00:00Z","temperature":6,"vpd":7}}`, string(b))
}

func TestSensorsHandlerOrigin(t *testing.T) {
	s := newTestServer(nil)
	req := httptest.NewRequest(http.MethodGet, "/sensors", nil)
	req.Header.Set("Origin", "ORIGIN")
	w := httptest.NewRecorder()
	s.sensorsHandler(w, req)
	res := w.Result()

	assert.Equal(t, http.StatusOK, res.StatusCode)
	assert.Equal(t, res.Header.Get("Access-Control-Allow-Origin"), "ORIGIN")
}

func TestBackgroundUpdater(t *testing.T) {
	client := &mockSPClient{}
	s := newTestServer(client)

	// kill the server after 700 milliseconds...
	go func() {
		time.Sleep(700 * time.Millisecond)
		s.appCancel()
	}()

	s.ListenAndServe()

	// in 700 milliseconds, we should have the initial LastSample call, and
	// then two more (one around 250ms, and one around 500ms)
	assert.Equal(t, 3, client.calls)
}
