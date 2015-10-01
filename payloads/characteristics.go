package payloads

import (
	"encoding/json"

	"github.com/thermokarst/bactdb/models"
)

type CharacteristicPayload struct {
	Characteristic *models.Characteristic     `json:"characteristic"`
	Measurements   *models.Measurements       `json:"measurements"`
	Strains        *models.Strains            `json:"strains"`
	Species        *models.ManySpecies        `json:"species"`
	Meta           *models.CharacteristicMeta `json:"meta"`
}

type CharacteristicsPayload struct {
	Characteristics *models.Characteristics    `json:"characteristics"`
	Measurements    *models.Measurements       `json:"measurements"`
	Strains         *models.Strains            `json:"strains"`
	Species         *models.ManySpecies        `json:"species"`
	Meta            *models.CharacteristicMeta `json:"meta"`
}

func (c *CharacteristicPayload) Marshal() ([]byte, error) {
	return json.Marshal(c)
}

func (c *CharacteristicsPayload) Marshal() ([]byte, error) {
	return json.Marshal(c)
}
