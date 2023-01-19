package sensorpush

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestSamples(t *testing.T) {
	client, calls := NewMockedClient(t)

	calls.
		On("Call",
			"samples",
			mock.AnythingOfType("*sensorpush.samplesRequest"),
			mock.AnythingOfType("*sensorpush.samplesResponse")).
		Run(func(args mock.Arguments) {
			r := args.Get(2).(*samplesResponse)
			(*r).Sensors = make(map[string][]Sample)
			(*r).Sensors["DEVICE-ID"] = []Sample{{Temperature: 70.0}}
		}).
		Return(nil)

	s, err := client.LastSample("DEVICE-ID")
	assert.NoError(t, err)
	assert.NotNil(t, s)
	assert.Equal(t, float32(70.0), s.Temperature)
}

func TestSamplesError(t *testing.T) {
	client, calls := NewMockedClient(t)

	anErr := errors.New("some error")

	calls.
		On("Call",
			"samples",
			mock.AnythingOfType("*sensorpush.samplesRequest"),
			mock.AnythingOfType("*sensorpush.samplesResponse")).
		Return(anErr)

	s, err := client.LastSample("DEVICE-ID")
	assert.Equal(t, anErr, err)
	assert.Nil(t, s)
}
