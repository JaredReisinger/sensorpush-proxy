package sensorpush

import "log/slog"

type sensorsRequest struct{}

// Sensor is a single sensor descriptor.
type Sensor struct {
	Active  bool   `json:"active"`
	Address string `json:"address"`
	// Alerts         struct{} `json:"alerts"`
	BatteryVoltage float32 `json:"battery_voltage"`
	// Calibration
	DeviceID string `json:"deviceId"`
	ID       string `json:"id"`
	Name     string `json:"name"`
	// Rssi     float32 `json:"rssi"`
	// Tags
	Type string `json:"type"`
}

// Sensors is a map of sensor descriptors, keyed by sensor device ID.
type Sensors map[string]Sensor

// Sensors gets a list of sensors
func (c *Client) Sensors() (*Sensors, error) {
	sensors := &Sensors{}
	err := c.authCall("devices/sensors", &sensorsRequest{}, sensors, 0)
	if err != nil {
		slog.Warn("unable to get sensors", "error", err)
		return nil, err
	}

	// slog.Debug("got sensors", "sensors", sensors)
	return sensors, nil
}
