package models

import (
	"errors"
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

type MockTextMeasurementTypesService struct {
	Get_ func(id int64) (*TextMeasurementType, error)
}

var _ TextMeasurementTypesService = &MockTextMeasurementTypesService{}

func (s *MockTextMeasurementTypesService) Get(id int64) (*TextMeasurementType, error) {
	if s.Get_ == nil {
		return nil, nil
	}
	return s.Get_(id)
}
