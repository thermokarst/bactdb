package models

import (
	"errors"
	"strconv"
	"time"

	"github.com/lib/pq"
	"github.com/thermokarst/bactdb/router"
)

// An Observation Type is a lookup type
type ObservationType struct {
	Id                  int64       `json:"id,omitempty"`
	ObservationTypeName string      `db:"observation_type_name" json:"observation_type_name"`
	CreatedAt           time.Time   `db:"created_at" json:"created_at"`
	UpdatedAt           time.Time   `db:"updated_at" json:"updated_at"`
	DeletedAt           pq.NullTime `db:"deleted_at" json:"deleted_at"`
}

func NewObservationType() *ObservationType {
	return &ObservationType{
		ObservationTypeName: "Test Obs Type",
	}
}

type ObservationTypesService interface {
	// Get an observation type
	Get(id int64) (*ObservationType, error)
}

var (
	ErrObservationTypeNotFound = errors.New("observation type not found")
)

type observationTypesService struct {
	client *Client
}

func (s *observationTypesService) Get(id int64) (*ObservationType, error) {
	strId := strconv.FormatInt(id, 10)

	url, err := s.client.url(router.ObservationType, map[string]string{"Id": strId}, nil)
	if err != nil {
		return nil, err
	}

	req, err := s.client.NewRequest("GET", url.String(), nil)
	if err != nil {
		return nil, err
	}

	var observation_type *ObservationType
	_, err = s.client.Do(req, &observation_type)
	if err != nil {
		return nil, err
	}

	return observation_type, nil
}

type MockObservationTypesService struct {
	Get_ func(id int64) (*ObservationType, error)
}

var _ ObservationTypesService = &MockObservationTypesService{}

func (s *MockObservationTypesService) Get(id int64) (*ObservationType, error) {
	if s.Get_ == nil {
		return nil, nil
	}
	return s.Get_(id)
}
