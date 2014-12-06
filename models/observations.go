package models

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/thermokarst/bactdb/router"
)

// An Observation is a lookup type
type Observation struct {
	Id                int64     `json:"id,omitempty"`
	ObservationName   string    `db:"observation_name" json:"observationName"`
	ObservationTypeId int64     `db:"observation_type_id" json:"observationTypeId"`
	CreatedAt         time.Time `db:"created_at" json:"createdAt"`
	UpdatedAt         time.Time `db:"updated_at" json:"updatedAt"`
	DeletedAt         NullTime  `db:"deleted_at" json:"deletedAt"`
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

	// Update an observation
	Update(id int64, Observation *Observation) (updated bool, err error)

	// Delete an observation
	Delete(id int64) (deleted bool, err error)
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

func (s *observationsService) Update(id int64, observation *Observation) (bool, error) {
	strId := strconv.FormatInt(id, 10)

	url, err := s.client.url(router.UpdateObservation, map[string]string{"Id": strId}, nil)
	if err != nil {
		return false, err
	}

	req, err := s.client.NewRequest("PUT", url.String(), observation)
	if err != nil {
		return false, err
	}

	resp, err := s.client.Do(req, &observation)
	if err != nil {
		return false, err
	}

	return resp.StatusCode == http.StatusOK, nil
}

func (s *observationsService) Delete(id int64) (bool, error) {
	strId := strconv.FormatInt(id, 10)

	url, err := s.client.url(router.DeleteObservation, map[string]string{"Id": strId}, nil)
	if err != nil {
		return false, err
	}

	req, err := s.client.NewRequest("DELETE", url.String(), nil)
	if err != nil {
		return false, err
	}

	var observation *Observation
	resp, err := s.client.Do(req, &observation)
	if err != nil {
		return false, err
	}

	return resp.StatusCode == http.StatusOK, nil
}

type MockObservationsService struct {
	Get_    func(id int64) (*Observation, error)
	List_   func(opt *ObservationListOptions) ([]*Observation, error)
	Create_ func(observation *Observation) (bool, error)
	Update_ func(id int64, observation *Observation) (bool, error)
	Delete_ func(id int64) (bool, error)
}

var _ ObservationsService = &MockObservationsService{}

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

func (s *MockObservationsService) Update(id int64, observation *Observation) (bool, error) {
	if s.Update_ == nil {
		return false, nil
	}
	return s.Update_(id, observation)
}

func (s *MockObservationsService) Delete(id int64) (bool, error) {
	if s.Delete_ == nil {
		return false, nil
	}
	return s.Delete_(id)
}
