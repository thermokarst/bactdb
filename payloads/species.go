package payloads

import (
	"encoding/json"

	"github.com/thermokarst/bactdb/models"
)

// Species is a payload that sideloads all of the necessary entities for a
// particular species.
type Species struct {
	Species *models.Species `json:"species"`
	Strains *models.Strains `json:"strains"`
}

// ManySpecies is a payload that sideloads all of the necessary entities for
// multiple species.
type ManySpecies struct {
	Species *models.ManySpecies `json:"species"`
	Strains *models.Strains     `json:"strains"`
}

// Marshal satisfies the CRUD interfaces.
func (s *Species) Marshal() ([]byte, error) {
	return json.Marshal(s)
}

// Marshal satisfies the CRUD interfaces.
func (s *ManySpecies) Marshal() ([]byte, error) {
	return json.Marshal(s)
}
