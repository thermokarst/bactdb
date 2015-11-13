package payloads

import (
	"encoding/json"

	"github.com/thermokarst/bactdb/models"
)

// Characteristic is a payload that sideloads all of the necessary entities for
// a particular characteristic.
type Characteristic struct {
	Characteristic *models.Characteristic `json:"characteristic"`
	Measurements   *models.Measurements   `json:"measurements"`
	Strains        *models.Strains        `json:"strains"`
	Species        *models.ManySpecies    `json:"species"`
}

// Characteristics is a payload that sideloads all of the necessary entities for
// multiple characteristics.
type Characteristics struct {
	Characteristics *models.Characteristics `json:"characteristics"`
	Measurements    *models.Measurements    `json:"measurements"`
	Strains         *models.Strains         `json:"strains"`
	Species         *models.ManySpecies     `json:"species"`
}

// Marshal satisfies the CRUD interfaces.
func (c *Characteristic) Marshal() ([]byte, error) {
	return json.Marshal(c)
}

// Marshal satisfies the CRUD interfaces.
func (c *Characteristics) Marshal() ([]byte, error) {
	return json.Marshal(c)
}
