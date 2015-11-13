package payloads

import (
	"encoding/json"

	"github.com/thermokarst/bactdb/models"
)

// Strain is a payload that sideloads all of the necessary entities for a
// particular strain.
type Strain struct {
	Strain          *models.Strain          `json:"strain"`
	Species         *models.ManySpecies     `json:"species"`
	Characteristics *models.Characteristics `json:"characteristics"`
	Measurements    *models.Measurements    `json:"measurements"`
}

// Strains is a payload that sideloads all of the necessary entities for
// multiple strains.
type Strains struct {
	Strains         *models.Strains         `json:"strains"`
	Species         *models.ManySpecies     `json:"species"`
	Characteristics *models.Characteristics `json:"characteristics"`
	Measurements    *models.Measurements    `json:"measurements"`
}

// Marshal satisfies the CRUD interfaces.
func (s *Strain) Marshal() ([]byte, error) {
	return json.Marshal(s)
}

// Marshal satisfies the CRUD interfaces.
func (s *Strains) Marshal() ([]byte, error) {
	return json.Marshal(s)
}
