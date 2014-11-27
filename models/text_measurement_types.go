package models

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/lib/pq"
	"github.com/thermokarst/bactdb/router"
)

// A TextMeasurementType is a lookup type
type TextMeasurementType struct {
	Id                  int64       `json:"id,omitempty"`
	TextMeasurementName string      `db:"text_measurement_name" json:"textMeasurementName"`
	CreatedAt           time.Time   `db:"created_at" json:"createdAt"`
	UpdatedAt           time.Time   `db:"updated_at" json:"updatedAt"`
	DeletedAt           pq.NullTime `db:"deleted_at" json:"deletedAt"`
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

type MockTextMeasurementTypesService struct {
	Get_    func(id int64) (*TextMeasurementType, error)
	List_   func(opt *TextMeasurementTypeListOptions) ([]*TextMeasurementType, error)
	Create_ func(text_measurement_type *TextMeasurementType) (bool, error)
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
