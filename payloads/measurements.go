package payloads

import (
	"encoding/json"

	"github.com/thermokarst/bactdb/models"
)

type MeasurementPayload struct {
	Measurement *models.Measurement `json:"measurement"`
}

type MeasurementsPayload struct {
	Strains         *models.Strains         `json:"strains"`
	Characteristics *models.Characteristics `json:"characteristics"`
	Measurements    *models.Measurements    `json:"measurements"`
}

func (m *MeasurementPayload) Marshal() ([]byte, error) {
	return json.Marshal(m)
}

func (m *MeasurementsPayload) Marshal() ([]byte, error) {
	return json.Marshal(m)
}
