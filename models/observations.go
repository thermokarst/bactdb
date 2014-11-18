package models

import (
	"errors"
	"strconv"
	"time"

	"github.com/lib/pq"
	"github.com/thermokarst/bactdb/router"
)

// An Observation is a lookup type
type Observation struct {
	Id                int64       `json:"id,omitempty"`
	ObservationName   string      `db:"observation_name" json:"observation_name"`
	ObservationTypeId int64       `db:"observation_type_id" json:"observation_type_id"`
	CreatedAt         time.Time   `db:"created_at" json:"created_at"`
	UpdatedAt         time.Time   `db:"updated_at" json:"updated_at"`
	DeletedAt         pq.NullTime `db:"deleted_at" json:"deleted_at"`
}

func NewObservation() *Observation {
	return &Observation{
		ObservationName: "Test Observation",
	}
}

type ObservationsService interface {
	// Get an observation
	Get(id int64) (*Observation, error)
}

var (
	ErrObservationNotFound = errors.New("observation not found")
)

type observationsService struct {
	client *Client
}

func (s *observationsService) Get(id int64) (*Observation, error) {
	strId := strconv.FormatInt(id, 10)

	url, err := s.client.url(router.Observation, map[string]string{"Id": strId}, nil)
	if err != nil {
		return nil, err
	}

	req, err := s.client.NewRequest("GET", url.String(), nil)
	if err != nil {
		return nil, err
	}

	var observation *Observation
	_, err = s.client.Do(req, &observation)
	if err != nil {
		return nil, err
	}

	return observation, nil
}

type MockObservationsService struct {
	Get_ func(id int64) (*Observation, error)
}

var _ObservationsService = &MockObservationsService{}

func (s *MockObservationsService) Get(id int64) (*Observation, error) {
	if s.Get_ == nil {
		return nil, nil
	}
	return s.Get_(id)
}
