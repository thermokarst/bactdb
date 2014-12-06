package models

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/thermokarst/bactdb/router"
)

// An Observation Type is a lookup type
type ObservationType struct {
	Id                  int64     `json:"id,omitempty"`
	ObservationTypeName string    `db:"observation_type_name" json:"observationTypeName"`
	CreatedAt           time.Time `db:"created_at" json:"createdAt"`
	UpdatedAt           time.Time `db:"updated_at" json:"updatedAt"`
	DeletedAt           NullTime  `db:"deleted_at" json:"deletedAt"`
}

func NewObservationType() *ObservationType {
	return &ObservationType{
		ObservationTypeName: "Test Obs Type",
	}
}

type ObservationTypesService interface {
	// Get an observation type
	Get(id int64) (*ObservationType, error)

	// List all observation types
	List(opt *ObservationTypeListOptions) ([]*ObservationType, error)

	// Create an observation type record
	Create(observation_type *ObservationType) (bool, error)

	// Update an existing observation type
	Update(id int64, observation_type *ObservationType) (updated bool, err error)

	// Delete an existing observation type
	Delete(id int64) (deleted bool, err error)
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

func (s *observationTypesService) Create(observation_type *ObservationType) (bool, error) {
	url, err := s.client.url(router.CreateObservationType, nil, nil)
	if err != nil {
		return false, err
	}

	req, err := s.client.NewRequest("POST", url.String(), observation_type)
	if err != nil {
		return false, err
	}

	resp, err := s.client.Do(req, &observation_type)
	if err != nil {
		return false, err
	}

	return resp.StatusCode == http.StatusCreated, nil
}

type ObservationTypeListOptions struct {
	ListOptions
}

func (s *observationTypesService) List(opt *ObservationTypeListOptions) ([]*ObservationType, error) {
	url, err := s.client.url(router.ObservationTypes, nil, opt)
	if err != nil {
		return nil, err
	}

	req, err := s.client.NewRequest("GET", url.String(), nil)
	if err != nil {
		return nil, err
	}

	var observation_types []*ObservationType
	_, err = s.client.Do(req, &observation_types)
	if err != nil {
		return nil, err
	}

	return observation_types, nil
}

func (s *observationTypesService) Update(id int64, observation_type *ObservationType) (bool, error) {
	strId := strconv.FormatInt(id, 10)

	url, err := s.client.url(router.UpdateObservationType, map[string]string{"Id": strId}, nil)
	if err != nil {
		return false, err
	}

	req, err := s.client.NewRequest("PUT", url.String(), observation_type)
	if err != nil {
		return false, err
	}

	resp, err := s.client.Do(req, &observation_type)
	if err != nil {
		return false, err
	}

	return resp.StatusCode == http.StatusOK, nil
}

func (s *observationTypesService) Delete(id int64) (bool, error) {
	strId := strconv.FormatInt(id, 10)

	url, err := s.client.url(router.DeleteObservationType, map[string]string{"Id": strId}, nil)
	if err != nil {
		return false, err
	}

	req, err := s.client.NewRequest("DELETE", url.String(), nil)
	if err != nil {
		return false, err
	}

	var observation_type *ObservationType
	resp, err := s.client.Do(req, &observation_type)
	if err != nil {
		return false, err
	}

	return resp.StatusCode == http.StatusOK, nil
}

type MockObservationTypesService struct {
	Get_    func(id int64) (*ObservationType, error)
	List_   func(opt *ObservationTypeListOptions) ([]*ObservationType, error)
	Create_ func(observation_type *ObservationType) (bool, error)
	Update_ func(id int64, observation_type *ObservationType) (bool, error)
	Delete_ func(id int64) (bool, error)
}

var _ ObservationTypesService = &MockObservationTypesService{}

func (s *MockObservationTypesService) Get(id int64) (*ObservationType, error) {
	if s.Get_ == nil {
		return nil, nil
	}
	return s.Get_(id)
}

func (s *MockObservationTypesService) Create(observation_type *ObservationType) (bool, error) {
	if s.Create_ == nil {
		return false, nil
	}
	return s.Create_(observation_type)
}

func (s *MockObservationTypesService) List(opt *ObservationTypeListOptions) ([]*ObservationType, error) {
	if s.List_ == nil {
		return nil, nil
	}
	return s.List_(opt)
}

func (s *MockObservationTypesService) Update(id int64, observation_type *ObservationType) (bool, error) {
	if s.Update_ == nil {
		return false, nil
	}
	return s.Update_(id, observation_type)
}

func (s *MockObservationTypesService) Delete(id int64) (bool, error) {
	if s.Delete_ == nil {
		return false, nil
	}
	return s.Delete_(id)
}
