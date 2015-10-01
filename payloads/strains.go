package payloads

import (
	"encoding/json"

	"github.com/thermokarst/bactdb/models"
)

type StrainPayload struct {
	Strain          *models.Strain          `json:"strain"`
	Species         *models.ManySpecies     `json:"species"`
	Characteristics *models.Characteristics `json:"characteristics"`
	Measurements    *models.Measurements    `json:"measurements"`
	Meta            *models.StrainMeta      `json:"meta"`
}

type StrainsPayload struct {
	Strains         *models.Strains         `json:"strains"`
	Species         *models.ManySpecies     `json:"species"`
	Characteristics *models.Characteristics `json:"characteristics"`
	Measurements    *models.Measurements    `json:"measurements"`
	Meta            *models.StrainMeta      `json:"meta"`
}

func (s *StrainPayload) Marshal() ([]byte, error) {
	return json.Marshal(s)
}

func (s *StrainsPayload) Marshal() ([]byte, error) {
	return json.Marshal(s)
}
