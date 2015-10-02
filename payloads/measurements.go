package payloads

import (
	"encoding/json"

	"github.com/thermokarst/bactdb/models"
)

// Measurement is a payload that sideloads all of the necessary entities for
// a particular measurement.
type Measurement struct {
	Measurement *models.Measurement `json:"measurement"`
}

// Measurements is a payload that sideloads all of the necessary entities for
// multiple measurements.
type Measurements struct {
	Strains         *models.Strains         `json:"strains"`
	Characteristics *models.Characteristics `json:"characteristics"`
	Measurements    *models.Measurements    `json:"measurements"`
}

// Marshal satisfies the CRUD interfaces.
func (m *Measurement) Marshal() ([]byte, error) {
	return json.Marshal(m)
}

// Marshal satisfies the CRUD interfaces.
func (m *Measurements) Marshal() ([]byte, error) {
	return json.Marshal(m)
}
