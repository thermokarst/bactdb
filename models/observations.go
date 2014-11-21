package models

import (
	"errors"
	"net/http"
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

	// List all observations
	List(opt *ObservationListOptions) ([]*Observation, error)

	// Create an observation
	Create(observation *Observation) (bool, error)
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

func (s *observationsService) Create(observation *Observation) (bool, error) {
	url, err := s.client.url(router.CreateObservation, nil, nil)
	if err != nil {
		return false, err
	}

	req, err := s.client.NewRequest("POST", url.String(), observation)
	if err != nil {
		return false, err
	}

	resp, err := s.client.Do(req, &observation)
	if err != nil {
		return false, err
	}

	return resp.StatusCode == http.StatusCreated, nil
}

type ObservationListOptions struct {
	ListOptions
}

func (s *observationsService) List(opt *ObservationListOptions) ([]*Observation, error) {
	url, err := s.client.url(router.Observations, nil, opt)
	if err != nil {
		return nil, err
	}

	req, err := s.client.NewRequest("GET", url.String(), nil)
	if err != nil {
		return nil, err
	}

	var observations []*Observation
	_, err = s.client.Do(req, &observations)
	if err != nil {
		return nil, err
	}

	return observations, nil
}

type MockObservationsService struct {
	Get_    func(id int64) (*Observation, error)
	List_   func(opt *ObservationListOptions) ([]*Observation, error)
	Create_ func(observation *Observation) (bool, error)
}

var _ObservationsService = &MockObservationsService{}

func (s *MockObservationsService) Get(id int64) (*Observation, error) {
	if s.Get_ == nil {
		return nil, nil
	}
	return s.Get_(id)
}

func (s *MockObservationsService) Create(observation *Observation) (bool, error) {
	if s.Create_ == nil {
		return false, nil
	}
	return s.Create_(observation)
}

func (s *MockObservationsService) List(opt *ObservationListOptions) ([]*Observation, error) {
	if s.List_ == nil {
		return nil, nil
	}
	return s.List_(opt)
}
