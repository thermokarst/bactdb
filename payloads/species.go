package payloads

import (
	"encoding/json"

	"github.com/thermokarst/bactdb/models"
)

type Species struct {
	Species *models.Species     `json:"species"`
	Strains *models.Strains     `json:"strains"`
	Meta    *models.SpeciesMeta `json:"meta"`
}

type ManySpecies struct {
	Species *models.ManySpecies `json:"species"`
	Strains *models.Strains     `json:"strains"`
	Meta    *models.SpeciesMeta `json:"meta"`
}

func (s *Species) Marshal() ([]byte, error) {
	return json.Marshal(s)
}

func (s *ManySpecies) Marshal() ([]byte, error) {
	return json.Marshal(s)
}
