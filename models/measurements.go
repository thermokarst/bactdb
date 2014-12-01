package models

import (
	"database/sql"
	"errors"
	"strconv"
	"time"

	"github.com/lib/pq"
	"github.com/thermokarst/bactdb/router"
)

// A Measurement is the main data type for this application
// There are two types of supported measurements: text & numerical. The table
// has a constraint that will allow one or the other for a particular
// combination of strain & observation, but not both.
type Measurement struct {
	Id                    int64           `json:"id,omitempty"`
	StrainId              int64           `db:"strain_id" json:"strainId"`
	ObservationId         int64           `db:"observation_id" json:"observationId"`
	TextMeasurementTypeId sql.NullInt64   `db:"text_measurement_type_id" json:"textMeasurementTypeId"`
	MeasurementValue      sql.NullFloat64 `db:"measurement_value" json:"measurementValue"`
	ConfidenceInterval    sql.NullFloat64 `db:"confidence_interval" json:"confidenceInterval"`
	UnitTypeId            sql.NullInt64   `db:"unit_type_id" json:"unitTypeId"`
	CreatedAt             time.Time       `db:"created_at" json:"createdAt"`
	UpdatedAt             time.Time       `db:"updated_at" json:"updatedAt"`
	DeletedAt             pq.NullTime     `db:"deleted_at" json:"deletedAt"`
}

func NewMeasurement() *Measurement {
	return &Measurement{
		MeasurementValue: sql.NullFloat64{Float64: 1.23, Valid: true},
	}
}

type MeasurementsService interface {
	// Get a measurement
	Get(id int64) (*Measurement, error)
}

var (
	ErrMeasurementNotFound = errors.New("measurement not found")
)

type measurementsService struct {
	client *Client
}

func (s *measurementsService) Get(id int64) (*Measurement, error) {
	strId := strconv.FormatInt(id, 10)

	url, err := s.client.url(router.Measurement, map[string]string{"Id": strId}, nil)
	if err != nil {
		return nil, err
	}

	req, err := s.client.NewRequest("GET", url.String(), nil)
	if err != nil {
		return nil, err
	}

	var measurement *Measurement
	_, err = s.client.Do(req, &measurement)
	if err != nil {
		return nil, err
	}

	return measurement, nil
}

type MockMeasurementsService struct {
	Get_ func(id int64) (*Measurement, error)
}

var _ MeasurementsService = &MockMeasurementsService{}

func (s *MockMeasurementsService) Get(id int64) (*Measurement, error) {
	if s.Get_ == nil {
		return nil, nil
	}
	return s.Get_(id)
}
