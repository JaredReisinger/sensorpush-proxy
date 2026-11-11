package core

import (
	"testing"

	"github.com/jaredreisinger/sensorpush-proxy/pkg/sensorpush"
	"github.com/stretchr/testify/assert"
)

type mockQueryClient struct {
	gateways *sensorpush.Gateways
	sensors  *sensorpush.Sensors
}

func (c *mockQueryClient) Gateways() (*sensorpush.Gateways, error) {
	return c.gateways, nil
}

func (c *mockQueryClient) Sensors() (*sensorpush.Sensors, error) {
	return c.sensors, nil
}

func TestQuery(t *testing.T) {
	// This really isn't testing anything other than that it doesn't error on
	// data. We need proper end-service mocking to test the whole flow.
	err := query(&mockQueryClient{})
	assert.Nil(t, err)

	err = query(&mockQueryClient{
		gateways: &sensorpush.Gateways{
			"GW": {},
		},
		sensors: &sensorpush.Sensors{
			"SN": {},
		},
	})
	assert.Nil(t, err)
}
