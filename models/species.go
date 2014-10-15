package models

import (
	"errors"
	"strconv"
	"time"

	"github.com/thermokarst/bactdb/router"
)

// A Species is a high-level classifier in bactdb.
type Species struct {
	Id          int64     `json:"id,omitempty"`
	GenusId     int64     `json:"genus_id"`
	SpeciesName string    `json:"species_name"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	DeletedAt   time.Time `json:"deleted_at"`
}

// SpeciesService interacts with the species-related endpoints in bactdb's API.
type SpeciesService interface {
	// Get a species
	Get(id int64) (*Species, error)
}

var (
	ErrSpeciesNotFound = errors.New("species not found")
)

type speciesService struct {
	client *Client
}

func (s *speciesService) Get(id int64) (*Species, error) {
	// Pass in key value pairs as strings, sp that the gorilla mux URL generation is happy
	strId := strconv.FormatInt(id, 10)

	url, err := s.client.url(router.Species, map[string]string{"Id": strId}, nil)
	if err != nil {
		return nil, err
	}

	req, err := s.client.NewRequest("GET", url.String(), nil)
	if err != nil {
		return nil, err
	}

	var species *Species
	_, err = s.client.Do(req, &species)
	if err != nil {
		return nil, err
	}

	return species, nil
}

type MockSpeciesService struct {
	Get_ func(id int64) (*Species, error)
}

var _ SpeciesService = &MockSpeciesService{}

func (s *MockSpeciesService) Get(id int64) (*Species, error) {
	if s.Get_ == nil {
		return nil, nil
	}
	return s.Get_(id)
}
