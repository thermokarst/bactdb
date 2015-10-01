package payloads

import (
	"encoding/json"

	"github.com/thermokarst/bactdb/models"
)

type SpeciesPayload struct {
	Species *models.Species     `json:"species"`
	Strains *models.Strains     `json:"strains"`
	Meta    *models.SpeciesMeta `json:"meta"`
}

type ManySpeciesPayload struct {
	Species *models.ManySpecies `json:"species"`
	Strains *models.Strains     `json:"strains"`
	Meta    *models.SpeciesMeta `json:"meta"`
}

func (s *SpeciesPayload) Marshal() ([]byte, error) {
	return json.Marshal(s)
}

func (s *ManySpeciesPayload) Marshal() ([]byte, error) {
	return json.Marshal(s)
}
