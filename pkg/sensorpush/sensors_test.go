package sensorpush

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestSensors(t *testing.T) {
	client, calls := NewMockedClient(t)

	calls.
		On("Call",
			"devices/sensors",
			mock.AnythingOfType("*sensorpush.sensorsRequest"),
			mock.AnythingOfType("*sensorpush.Sensors")).
		Run(func(args mock.Arguments) {
			r := args.Get(2).(*Sensors)
			(*r)["key"] = Sensor{DeviceID: "DEVICE-ID"}
		}).
		Return(nil)

	s, err := client.Sensors()
	assert.NoError(t, err)
	assert.NotNil(t, s)
	assert.Equal(t, "DEVICE-ID", (*s)["key"].DeviceID)
}

func TestSensorsError(t *testing.T) {
	client, calls := NewMockedClient(t)

	anErr := errors.New("some error")

	calls.
		On("Call",
			"devices/sensors",
			mock.AnythingOfType("*sensorpush.sensorsRequest"),
			mock.AnythingOfType("*sensorpush.Sensors")).
		Return(anErr)

	s, err := client.Sensors()
	assert.Equal(t, anErr, err)
	assert.Nil(t, s)
}
