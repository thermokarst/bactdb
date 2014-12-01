package models

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/lib/pq"
	"github.com/thermokarst/bactdb/router"
)

// A UnitType is a lookup type
type UnitType struct {
	Id        int64       `json:"id,omitempty"`
	Name      string      `db:"name" json:"name"`
	Symbol    string      `db:"symbol" json:"symbol"`
	CreatedAt time.Time   `db:"created_at" json:"createdAt"`
	UpdatedAt time.Time   `db:"updated_at" json:"updatedAt"`
	DeletedAt pq.NullTime `db:"deleted_at" json:"deletedAt"`
}

func NewUnitType() *UnitType {
	return &UnitType{
		Name:   "Test Unit Type",
		Symbol: "x",
	}
}

type UnitTypesService interface {
	// Get a unit type
	Get(id int64) (*UnitType, error)

	// List all unit types
	List(opt *UnitTypeListOptions) ([]*UnitType, error)

	// Create a unit type
	Create(unit_type *UnitType) (bool, error)

	// Update a unit type
	Update(id int64, UnitType *UnitType) (bool, error)

	// Delete a unit type
	Delete(id int64) (deleted bool, err error)
}

var (
	ErrUnitTypeNotFound = errors.New("unit type not found")
)

type unitTypesService struct {
	client *Client
}

func (s *unitTypesService) Get(id int64) (*UnitType, error) {
	strId := strconv.FormatInt(id, 10)

	url, err := s.client.url(router.UnitType, map[string]string{"Id": strId}, nil)
	if err != nil {
		return nil, err
	}

	req, err := s.client.NewRequest("GET", url.String(), nil)
	if err != nil {
		return nil, err
	}

	var unit_type *UnitType
	_, err = s.client.Do(req, &unit_type)
	if err != nil {
		return nil, err
	}

	return unit_type, nil
}

func (s *unitTypesService) Create(unit_type *UnitType) (bool, error) {
	url, err := s.client.url(router.CreateUnitType, nil, nil)
	if err != nil {
		return false, err
	}

	req, err := s.client.NewRequest("POST", url.String(), unit_type)
	if err != nil {
		return false, err
	}

	resp, err := s.client.Do(req, &unit_type)
	if err != nil {
		return false, err
	}

	return resp.StatusCode == http.StatusCreated, nil
}

type UnitTypeListOptions struct {
	ListOptions
}

func (s *unitTypesService) List(opt *UnitTypeListOptions) ([]*UnitType, error) {
	url, err := s.client.url(router.UnitTypes, nil, opt)
	if err != nil {
		return nil, err
	}

	req, err := s.client.NewRequest("GET", url.String(), nil)
	if err != nil {
		return nil, err
	}

	var unit_types []*UnitType
	_, err = s.client.Do(req, &unit_types)
	if err != nil {
		return nil, err
	}

	return unit_types, nil
}

func (s *unitTypesService) Update(id int64, unit_type *UnitType) (bool, error) {
	strId := strconv.FormatInt(id, 10)

	url, err := s.client.url(router.UpdateUnitType, map[string]string{"Id": strId}, nil)
	if err != nil {
		return false, err
	}

	req, err := s.client.NewRequest("PUT", url.String(), unit_type)
	if err != nil {
		return false, err
	}

	resp, err := s.client.Do(req, &unit_type)
	if err != nil {
		return false, err
	}

	return resp.StatusCode == http.StatusOK, nil
}

func (s *unitTypesService) Delete(id int64) (bool, error) {
	strId := strconv.FormatInt(id, 10)

	url, err := s.client.url(router.DeleteUnitType, map[string]string{"Id": strId}, nil)
	if err != nil {
		return false, err
	}

	req, err := s.client.NewRequest("DELETE", url.String(), nil)
	if err != nil {
		return false, err
	}

	var unit_type *UnitType
	resp, err := s.client.Do(req, &unit_type)
	if err != nil {
		return false, err
	}

	return resp.StatusCode == http.StatusOK, nil
}

type MockUnitTypesService struct {
	Get_    func(id int64) (*UnitType, error)
	List_   func(opt *UnitTypeListOptions) ([]*UnitType, error)
	Create_ func(unit_type *UnitType) (bool, error)
	Update_ func(id int64, unit_type *UnitType) (bool, error)
	Delete_ func(id int64) (bool, error)
}

var _ UnitTypesService = &MockUnitTypesService{}

func (s *MockUnitTypesService) Get(id int64) (*UnitType, error) {
	if s.Get_ == nil {
		return nil, nil
	}
	return s.Get_(id)
}

func (s *MockUnitTypesService) Create(unit_type *UnitType) (bool, error) {
	if s.Create_ == nil {
		return false, nil
	}
	return s.Create_(unit_type)
}

func (s *MockUnitTypesService) List(opt *UnitTypeListOptions) ([]*UnitType, error) {
	if s.List_ == nil {
		return nil, nil
	}
	return s.List_(opt)
}

func (s *MockUnitTypesService) Update(id int64, unit_type *UnitType) (bool, error) {
	if s.Update_ == nil {
		return false, nil
	}
	return s.Update_(id, unit_type)
}

func (s *MockUnitTypesService) Delete(id int64) (bool, error) {
	if s.Delete_ == nil {
		return false, nil
	}
	return s.Delete_(id)
}
