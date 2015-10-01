package payloads

import (
	"encoding/json"

	"github.com/thermokarst/bactdb/models"
)

type Measurement struct {
	Measurement *models.Measurement `json:"measurement"`
}

type Measurements struct {
	Strains         *models.Strains         `json:"strains"`
	Characteristics *models.Characteristics `json:"characteristics"`
	Measurements    *models.Measurements    `json:"measurements"`
}

func (m *Measurement) Marshal() ([]byte, error) {
	return json.Marshal(m)
}

func (m *Measurements) Marshal() ([]byte, error) {
	return json.Marshal(m)
}
