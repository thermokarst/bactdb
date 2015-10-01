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

type UserPayload struct {
	User *models.User     `json:"user"`
	Meta *models.UserMeta `json:"meta"`
}

func (u *UserPayload) Marshal() ([]byte, error) {
	return json.Marshal(u)
}
