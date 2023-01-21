package sensorpush

import (
	"log"
)

type gatewaysRequest struct{}

// Gateway is a single gateway descriptor.
type Gateway struct {
	LastAlert string `json:"last_alert"`
	LastSeen  string `json:"last_seen"`
	Message   string `json:"message"`
	Name      string `json:"name"`
	Paired    bool   `json:"paired"`
	// Samples
	// Tags
	Version string `json:"version"`
	ID      string `json:"id"`
}

// Gateways is a map of gateway descriptors, keyed by gateway ID.
type Gateways map[string]Gateway

// Gateways gets a list of gateways.
func (c *Client) Gateways() (*Gateways, error) {
	gateways := &Gateways{}
	err := c.authCall("devices/gateways", &gatewaysRequest{}, gateways, 0)
	if err != nil {
		log.Printf("unable to get gateways: %+v", err)
		return nil, err
	}

	// log.Printf("gateways: %+v", gateways)
	return gateways, nil
}
