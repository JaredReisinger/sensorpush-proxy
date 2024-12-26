package cmd

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/jaredreisinger/sensorpush-proxy/pkg/sensorpush"
	"github.com/stretchr/testify/assert"
)

func TestSensorsHandler(t *testing.T) {
	// Set up some bogus data...
	lastSamples["DUMMY"] = sensorpush.Sample{
		Altitude:           1,
		BarometricPressure: 2,
		Dewpoint:           3,
		Humidity:           4,
		Observed:           time.Date(2000, 1, 1, 12, 0, 0, 0, time.UTC),
		Temperature:        6,
		Vpd:                7,
	}

	req := httptest.NewRequest(http.MethodGet, "/sensors", nil)
	w := httptest.NewRecorder()
	sensorsHandler(w, req)
	res := w.Result()

	assert.Equal(t, http.StatusOK, res.StatusCode)
	assert.Equal(t, "application/json", res.Header.Get("Content-Type"))

	defer res.Body.Close()
	b, err := io.ReadAll(res.Body)
	assert.NoError(t, err)
	assert.Equal(t, `{"DUMMY":{"altitude":1,"barometric_pressure":2,"dewpoint":3,"humidity":4,"observed":"2000-01-01T12:00:00Z","temperature":6,"vpd":7}}`, string(b))
}

func TestSensorsHandlerOrigin(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/sensors", nil)
	req.Header.Set("Origin", "ORIGIN")
	w := httptest.NewRecorder()
	sensorsHandler(w, req)
	res := w.Result()

	assert.Equal(t, http.StatusOK, res.StatusCode)
	assert.Equal(t, res.Header.Get("Access-Control-Allow-Origin"), "ORIGIN")
}
