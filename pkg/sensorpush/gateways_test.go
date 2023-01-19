package sensorpush

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGateways(t *testing.T) {
	client, calls := NewMockedClient(t)

	calls.
		On("Call",
			"devices/gateways",
			mock.AnythingOfType("*sensorpush.gatewaysRequest"),
			mock.AnythingOfType("*sensorpush.Gateways")).
		Run(func(args mock.Arguments) {
			r := args.Get(2).(*Gateways)
			(*r)["key"] = Gateway{ID: "ID"}
		}).
		Return(nil)

	g, err := client.Gateways()
	assert.NoError(t, err)
	assert.NotNil(t, g)
	assert.Equal(t, "ID", (*g)["key"].ID)
}

func TestGatewaysError(t *testing.T) {
	client, calls := NewMockedClient(t)

	anErr := errors.New("some error")

	calls.
		On("Call",
			"devices/gateways",
			mock.AnythingOfType("*sensorpush.gatewaysRequest"),
			mock.AnythingOfType("*sensorpush.Gateways")).
		Return(anErr)

	g, err := client.Gateways()
	assert.Equal(t, anErr, err)
	assert.Nil(t, g)
}
