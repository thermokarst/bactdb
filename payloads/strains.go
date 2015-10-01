package payloads

import (
	"encoding/json"

	"github.com/thermokarst/bactdb/models"
)

type Strain struct {
	Strain          *models.Strain          `json:"strain"`
	Species         *models.ManySpecies     `json:"species"`
	Characteristics *models.Characteristics `json:"characteristics"`
	Measurements    *models.Measurements    `json:"measurements"`
	Meta            *models.StrainMeta      `json:"meta"`
}

type Strains struct {
	Strains         *models.Strains         `json:"strains"`
	Species         *models.ManySpecies     `json:"species"`
	Characteristics *models.Characteristics `json:"characteristics"`
	Measurements    *models.Measurements    `json:"measurements"`
	Meta            *models.StrainMeta      `json:"meta"`
}

func (s *Strain) Marshal() ([]byte, error) {
	return json.Marshal(s)
}

func (s *Strains) Marshal() ([]byte, error) {
	return json.Marshal(s)
}
