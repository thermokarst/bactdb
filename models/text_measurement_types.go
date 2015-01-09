package models

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/thermokarst/bactdb/router"
)

// A TextMeasurementType is a lookup type
type TextMeasurementType struct {
	Id                  int64     `json:"id,omitempty"`
	TextMeasurementName string    `db:"text_measurement_name" json:"textMeasurementName"`
	CreatedAt           time.Time `db:"created_at" json:"createdAt"`
	UpdatedAt           time.Time `db:"updated_at" json:"updatedAt"`
	DeletedAt           NullTime  `db:"deleted_at" json:"deletedAt"`
}

func (m *TextMeasurementType) String() string {
	return fmt.Sprintf("%v", *m)
}

func NewTextMeasurementType() *TextMeasurementType {
	return &TextMeasurementType{
		TextMeasurementName: "Test Text Measurement Type",
	}
}

type TextMeasurementTypesService interface {
	// Get a text measurement type
	Get(id int64) (*TextMeasurementType, error)

	// List all text measurement types
	List(opt *TextMeasurementTypeListOptions) ([]*TextMeasurementType, error)

	// Create a text measurement type
	Create(text_measurement_type *TextMeasurementType) (bool, error)

	// Update a text measurement type
	Update(id int64, TextMeasurementType *TextMeasurementType) (updated bool, err error)

	// Delete a text measurement type
	Delete(id int64) (deleted bool, err error)
}

var (
	ErrTextMeasurementTypeNotFound = errors.New("text measurement type not found")
)

type textMeasurementTypesService struct {
	client *Client
}

func (s *textMeasurementTypesService) Get(id int64) (*TextMeasurementType, error) {
	strId := strconv.FormatInt(id, 10)

	url, err := s.client.url(router.TextMeasurementType, map[string]string{"Id": strId}, nil)
	if err != nil {
		return nil, err
	}

	req, err := s.client.NewRequest("GET", url.String(), nil)
	if err != nil {
		return nil, err
	}

	var text_measurement_type *TextMeasurementType
	_, err = s.client.Do(req, &text_measurement_type)
	if err != nil {
		return nil, err
	}

	return text_measurement_type, nil
}

func (s *textMeasurementTypesService) Create(text_measurement_type *TextMeasurementType) (bool, error) {
	url, err := s.client.url(router.CreateTextMeasurementType, nil, nil)
	if err != nil {
		return false, err
	}

	req, err := s.client.NewRequest("POST", url.String(), text_measurement_type)
	if err != nil {
		return false, err
	}

	resp, err := s.client.Do(req, &text_measurement_type)
	if err != nil {
		return false, err
	}

	return resp.StatusCode == http.StatusCreated, nil
}

type TextMeasurementTypeListOptions struct {
	ListOptions
}

func (s *textMeasurementTypesService) List(opt *TextMeasurementTypeListOptions) ([]*TextMeasurementType, error) {
	url, err := s.client.url(router.TextMeasurementTypes, nil, opt)
	if err != nil {
		return nil, err
	}

	req, err := s.client.NewRequest("GET", url.String(), nil)
	if err != nil {
		return nil, err
	}

	var text_measurement_types []*TextMeasurementType
	_, err = s.client.Do(req, &text_measurement_types)
	if err != nil {
		return nil, err
	}

	return text_measurement_types, nil
}

func (s *textMeasurementTypesService) Update(id int64, text_measurement_type *TextMeasurementType) (bool, error) {
	strId := strconv.FormatInt(id, 10)

	url, err := s.client.url(router.UpdateTextMeasurementType, map[string]string{"Id": strId}, nil)
	if err != nil {
		return false, err
	}

	req, err := s.client.NewRequest("PUT", url.String(), text_measurement_type)
	if err != nil {
		return false, err
	}

	resp, err := s.client.Do(req, &text_measurement_type)
	if err != nil {
		return false, err
	}

	return resp.StatusCode == http.StatusOK, nil
}

func (s *textMeasurementTypesService) Delete(id int64) (bool, error) {
	strId := strconv.FormatInt(id, 10)

	url, err := s.client.url(router.DeleteTextMeasurementType, map[string]string{"Id": strId}, nil)
	if err != nil {
		return false, err
	}

	req, err := s.client.NewRequest("DELETE", url.String(), nil)
	if err != nil {
		return false, err
	}

	var text_measurement_type *TextMeasurementType
	resp, err := s.client.Do(req, &text_measurement_type)
	if err != nil {
		return false, err
	}

	return resp.StatusCode == http.StatusOK, nil
}

type MockTextMeasurementTypesService struct {
	Get_    func(id int64) (*TextMeasurementType, error)
	List_   func(opt *TextMeasurementTypeListOptions) ([]*TextMeasurementType, error)
	Create_ func(text_measurement_type *TextMeasurementType) (bool, error)
	Update_ func(id int64, text_measurement_type *TextMeasurementType) (bool, error)
	Delete_ func(id int64) (bool, error)
}

var _ TextMeasurementTypesService = &MockTextMeasurementTypesService{}

func (s *MockTextMeasurementTypesService) Get(id int64) (*TextMeasurementType, error) {
	if s.Get_ == nil {
		return nil, nil
	}
	return s.Get_(id)
}

func (s *MockTextMeasurementTypesService) Create(text_measurement_type *TextMeasurementType) (bool, error) {
	if s.Create_ == nil {
		return false, nil
	}
	return s.Create_(text_measurement_type)
}

func (s *MockTextMeasurementTypesService) List(opt *TextMeasurementTypeListOptions) ([]*TextMeasurementType, error) {
	if s.List_ == nil {
		return nil, nil
	}
	return s.List_(opt)
}

func (s *MockTextMeasurementTypesService) Update(id int64, text_measurement_type *TextMeasurementType) (bool, error) {
	if s.Update_ == nil {
		return false, nil
	}
	return s.Update_(id, text_measurement_type)
}

func (s *MockTextMeasurementTypesService) Delete(id int64) (bool, error) {
	if s.Delete_ == nil {
		return false, nil
	}
	return s.Delete_(id)
}
