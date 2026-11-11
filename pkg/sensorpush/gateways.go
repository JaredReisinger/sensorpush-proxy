package sensorpush

import "log/slog"

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
		slog.Warn("unable to get gateways", "error", err)
		return nil, err
	}

	// slog.Debug("got gateways", "gateways", gateways)
	return gateways, nil
}
