package models

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/thermokarst/bactdb/router"
)

// A Characteristic is a lookup type
type Characteristic struct {
	Id                   int64     `json:"id,omitempty"`
	CharacteristicName   string    `db:"characteristic_name" json:"characteristicName"`
	CharacteristicTypeId int64     `db:"characteristic_type_id" json:"characteristicTypeId"`
	CreatedAt            time.Time `db:"created_at" json:"createdAt"`
	UpdatedAt            time.Time `db:"updated_at" json:"updatedAt"`
	DeletedAt            NullTime  `db:"deleted_at" json:"deletedAt"`
}

func (m *Characteristic) String() string {
	return fmt.Sprintf("%v", *m)
}

func NewCharacteristic() *Characteristic {
	return &Characteristic{
		CharacteristicName: "Test Characteristic",
	}
}

type CharacteristicsService interface {
	// Get an characteristic
	Get(id int64) (*Characteristic, error)

	// List all characteristics
	List(opt *CharacteristicListOptions) ([]*Characteristic, error)

	// Create an characteristic
	Create(characteristic *Characteristic) (bool, error)

	// Update an characteristic
	Update(id int64, Characteristic *Characteristic) (updated bool, err error)

	// Delete an characteristic
	Delete(id int64) (deleted bool, err error)
}

var (
	ErrCharacteristicNotFound = errors.New("characteristic not found")
)

type characteristicsService struct {
	client *Client
}

func (s *characteristicsService) Get(id int64) (*Characteristic, error) {
	strId := strconv.FormatInt(id, 10)

	url, err := s.client.url(router.Characteristic, map[string]string{"Id": strId}, nil)
	if err != nil {
		return nil, err
	}

	req, err := s.client.NewRequest("GET", url.String(), nil)
	if err != nil {
		return nil, err
	}

	var characteristic *Characteristic
	_, err = s.client.Do(req, &characteristic)
	if err != nil {
		return nil, err
	}

	return characteristic, nil
}

func (s *characteristicsService) Create(characteristic *Characteristic) (bool, error) {
	url, err := s.client.url(router.CreateCharacteristic, nil, nil)
	if err != nil {
		return false, err
	}

	req, err := s.client.NewRequest("POST", url.String(), characteristic)
	if err != nil {
		return false, err
	}

	resp, err := s.client.Do(req, &characteristic)
	if err != nil {
		return false, err
	}

	return resp.StatusCode == http.StatusCreated, nil
}

type CharacteristicListOptions struct {
	ListOptions
}

func (s *characteristicsService) List(opt *CharacteristicListOptions) ([]*Characteristic, error) {
	url, err := s.client.url(router.Characteristics, nil, opt)
	if err != nil {
		return nil, err
	}

	req, err := s.client.NewRequest("GET", url.String(), nil)
	if err != nil {
		return nil, err
	}

	var characteristics []*Characteristic
	_, err = s.client.Do(req, &characteristics)
	if err != nil {
		return nil, err
	}

	return characteristics, nil
}

func (s *characteristicsService) Update(id int64, characteristic *Characteristic) (bool, error) {
	strId := strconv.FormatInt(id, 10)

	url, err := s.client.url(router.UpdateCharacteristic, map[string]string{"Id": strId}, nil)
	if err != nil {
		return false, err
	}

	req, err := s.client.NewRequest("PUT", url.String(), characteristic)
	if err != nil {
		return false, err
	}

	resp, err := s.client.Do(req, &characteristic)
	if err != nil {
		return false, err
	}

	return resp.StatusCode == http.StatusOK, nil
}

func (s *characteristicsService) Delete(id int64) (bool, error) {
	strId := strconv.FormatInt(id, 10)

	url, err := s.client.url(router.DeleteCharacteristic, map[string]string{"Id": strId}, nil)
	if err != nil {
		return false, err
	}

	req, err := s.client.NewRequest("DELETE", url.String(), nil)
	if err != nil {
		return false, err
	}

	var characteristic *Characteristic
	resp, err := s.client.Do(req, &characteristic)
	if err != nil {
		return false, err
	}

	return resp.StatusCode == http.StatusOK, nil
}

type MockCharacteristicsService struct {
	Get_    func(id int64) (*Characteristic, error)
	List_   func(opt *CharacteristicListOptions) ([]*Characteristic, error)
	Create_ func(characteristic *Characteristic) (bool, error)
	Update_ func(id int64, characteristic *Characteristic) (bool, error)
	Delete_ func(id int64) (bool, error)
}

var _ CharacteristicsService = &MockCharacteristicsService{}

func (s *MockCharacteristicsService) Get(id int64) (*Characteristic, error) {
	if s.Get_ == nil {
		return nil, nil
	}
	return s.Get_(id)
}

func (s *MockCharacteristicsService) Create(characteristic *Characteristic) (bool, error) {
	if s.Create_ == nil {
		return false, nil
	}
	return s.Create_(characteristic)
}

func (s *MockCharacteristicsService) List(opt *CharacteristicListOptions) ([]*Characteristic, error) {
	if s.List_ == nil {
		return nil, nil
	}
	return s.List_(opt)
}

func (s *MockCharacteristicsService) Update(id int64, characteristic *Characteristic) (bool, error) {
	if s.Update_ == nil {
		return false, nil
	}
	return s.Update_(id, characteristic)
}

func (s *MockCharacteristicsService) Delete(id int64) (bool, error) {
	if s.Delete_ == nil {
		return false, nil
	}
	return s.Delete_(id)
}
