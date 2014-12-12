package models

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/thermokarst/bactdb/router"
)

// A Characteristic Type is a lookup type
type CharacteristicType struct {
	Id                     int64     `json:"id,omitempty"`
	CharacteristicTypeName string    `db:"characteristic_type_name" json:"characteristicTypeName"`
	CreatedAt              time.Time `db:"created_at" json:"createdAt"`
	UpdatedAt              time.Time `db:"updated_at" json:"updatedAt"`
	DeletedAt              NullTime  `db:"deleted_at" json:"deletedAt"`
}

func NewCharacteristicType() *CharacteristicType {
	return &CharacteristicType{
		CharacteristicTypeName: "Test Char Type",
	}
}

type CharacteristicTypesService interface {
	// Get a characteristic type
	Get(id int64) (*CharacteristicType, error)

	// List all characteristic types
	List(opt *CharacteristicTypeListOptions) ([]*CharacteristicType, error)

	// Create a characteristic type record
	Create(characteristic_type *CharacteristicType) (bool, error)

	// Update an existing characteristic type
	Update(id int64, characteristic_type *CharacteristicType) (updated bool, err error)

	// Delete an existing characteristic type
	Delete(id int64) (deleted bool, err error)
}

var (
	ErrCharacteristicTypeNotFound = errors.New("characteristic type not found")
)

type characteristicTypesService struct {
	client *Client
}

func (s *characteristicTypesService) Get(id int64) (*CharacteristicType, error) {
	strId := strconv.FormatInt(id, 10)

	url, err := s.client.url(router.CharacteristicType, map[string]string{"Id": strId}, nil)
	if err != nil {
		return nil, err
	}

	req, err := s.client.NewRequest("GET", url.String(), nil)
	if err != nil {
		return nil, err
	}

	var characteristic_type *CharacteristicType
	_, err = s.client.Do(req, &characteristic_type)
	if err != nil {
		return nil, err
	}

	return characteristic_type, nil
}

func (s *characteristicTypesService) Create(characteristic_type *CharacteristicType) (bool, error) {
	url, err := s.client.url(router.CreateCharacteristicType, nil, nil)
	if err != nil {
		return false, err
	}

	req, err := s.client.NewRequest("POST", url.String(), characteristic_type)
	if err != nil {
		return false, err
	}

	resp, err := s.client.Do(req, &characteristic_type)
	if err != nil {
		return false, err
	}

	return resp.StatusCode == http.StatusCreated, nil
}

type CharacteristicTypeListOptions struct {
	ListOptions
}

func (s *characteristicTypesService) List(opt *CharacteristicTypeListOptions) ([]*CharacteristicType, error) {
	url, err := s.client.url(router.CharacteristicTypes, nil, opt)
	if err != nil {
		return nil, err
	}

	req, err := s.client.NewRequest("GET", url.String(), nil)
	if err != nil {
		return nil, err
	}

	var characteristic_types []*CharacteristicType
	_, err = s.client.Do(req, &characteristic_types)
	if err != nil {
		return nil, err
	}

	return characteristic_types, nil
}

func (s *characteristicTypesService) Update(id int64, characteristic_type *CharacteristicType) (bool, error) {
	strId := strconv.FormatInt(id, 10)

	url, err := s.client.url(router.UpdateCharacteristicType, map[string]string{"Id": strId}, nil)
	if err != nil {
		return false, err
	}

	req, err := s.client.NewRequest("PUT", url.String(), characteristic_type)
	if err != nil {
		return false, err
	}

	resp, err := s.client.Do(req, &characteristic_type)
	if err != nil {
		return false, err
	}

	return resp.StatusCode == http.StatusOK, nil
}

func (s *characteristicTypesService) Delete(id int64) (bool, error) {
	strId := strconv.FormatInt(id, 10)

	url, err := s.client.url(router.DeleteCharacteristicType, map[string]string{"Id": strId}, nil)
	if err != nil {
		return false, err
	}

	req, err := s.client.NewRequest("DELETE", url.String(), nil)
	if err != nil {
		return false, err
	}

	var characteristic_type *CharacteristicType
	resp, err := s.client.Do(req, &characteristic_type)
	if err != nil {
		return false, err
	}

	return resp.StatusCode == http.StatusOK, nil
}

type MockCharacteristicTypesService struct {
	Get_    func(id int64) (*CharacteristicType, error)
	List_   func(opt *CharacteristicTypeListOptions) ([]*CharacteristicType, error)
	Create_ func(characteristic_type *CharacteristicType) (bool, error)
	Update_ func(id int64, characteristic_type *CharacteristicType) (bool, error)
	Delete_ func(id int64) (bool, error)
}

var _ CharacteristicTypesService = &MockCharacteristicTypesService{}

func (s *MockCharacteristicTypesService) Get(id int64) (*CharacteristicType, error) {
	if s.Get_ == nil {
		return nil, nil
	}
	return s.Get_(id)
}

func (s *MockCharacteristicTypesService) Create(characteristic_type *CharacteristicType) (bool, error) {
	if s.Create_ == nil {
		return false, nil
	}
	return s.Create_(characteristic_type)
}

func (s *MockCharacteristicTypesService) List(opt *CharacteristicTypeListOptions) ([]*CharacteristicType, error) {
	if s.List_ == nil {
		return nil, nil
	}
	return s.List_(opt)
}

func (s *MockCharacteristicTypesService) Update(id int64, characteristic_type *CharacteristicType) (bool, error) {
	if s.Update_ == nil {
		return false, nil
	}
	return s.Update_(id, characteristic_type)
}

func (s *MockCharacteristicTypesService) Delete(id int64) (bool, error) {
	if s.Delete_ == nil {
		return false, nil
	}
	return s.Delete_(id)
}
