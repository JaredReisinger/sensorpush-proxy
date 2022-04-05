package sensorpush

import (
	"log"
	"time"
)

type samplesRequest struct {
	// Active    bool      `json:"active"`
	// Bulk      bool      `json:"bulk,omitempty"`
	Format    string     `json:"format,omitempty"`
	Limit     int        `json:"limit,omitempty"`
	Measures  []string   `json:"measures,omitempty"`
	Sensors   []string   `json:"sensors,omitempty"`
	StartTime *time.Time `json:"startTime,omitempty"`
	StopTime  *time.Time `json:"stopTime,omitempty"`
}

// Sample is a single sensor reading
type Sample struct {
	Altitude           float32   `json:"altitude"`
	BarometricPressure float32   `json:"barometric_pressure"`
	Dewpoint           float32   `json:"dewpoint"`
	Humidity           float32   `json:"humidity"`
	Observed           time.Time `json:"observed"`
	Temperature        float32   `json:"temperature"`
	Vpd                float32   `json:"vpd"`
}

type samplesResponse struct {
	LastTime     time.Time           `json:"last_time"`
	Sensors      map[string][]Sample `json:"sensors"`
	Status       string              `json:"status"`
	TotalSamples int                 `json:"total_samples"`
	TotalSensors int                 `json:"total_sensors"`
	Truncated    bool                `json:"truncated"`
}

// Sample gets the last sample for a given device
func (c *Client) Sample(deviceID string) (*Sample, error) {
	samples := &samplesResponse{}
	err := c.authCall("samples", &samplesRequest{Sensors: []string{deviceID}, Limit: 1}, samples, 0)
	if err != nil {
		log.Printf("unable to get samples: %+v", err)
		return nil, err
	}

	// log.Printf("got samples: %+v", samples)

	sample := samples.Sensors[deviceID][0]

	// log.Printf("SAMPLE: %s: %fF, %f%%RH", sample.Observed, sample.Temperature, sample.Humidity)

	return &sample, nil
}
