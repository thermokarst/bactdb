package models

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/thermokarst/bactdb/router"
)

// A Measurement is the main data type for this application
// There are two types of supported measurements: text & numerical. The table
// has a constraint that will allow one or the other for a particular
// combination of strain & characteristic, but not both.
type Measurement struct {
	Id                    int64       `json:"id,omitempty"`
	StrainId              int64       `db:"strain_id" json:"strain"`
	CharacteristicId      int64       `db:"characteristic_id" json:"characteristic"`
	TextMeasurementTypeId NullInt64   `db:"text_measurement_type_id" json:"textMeasurementTypeId"`
	TxtValue              NullString  `db:"txt_value" json:"txtValue"`
	NumValue              NullFloat64 `db:"num_value" json:"numValue"`
	ConfidenceInterval    NullFloat64 `db:"confidence_interval" json:"confidenceInterval"`
	UnitTypeId            NullInt64   `db:"unit_type_id" json:"unitTypeId"`
	Notes                 NullString  `db:"notes" json:"notes"`
	TestMethodId          NullInt64   `db:"test_method_id" json:"testMethodId"`
	CreatedAt             time.Time   `db:"created_at" json:"createdAt"`
	UpdatedAt             time.Time   `db:"updated_at" json:"updatedAt"`
}

type MeasurementJSON struct {
	Measurement *Measurement `json:"measurement"`
}

type MeasurementsJSON struct {
	Measurements []*Measurement `json:"measurements"`
}

func (m *Measurement) String() string {
	return fmt.Sprintf("%v", *m)
}

func NewMeasurement() *Measurement {
	return &Measurement{
		NumValue: NullFloat64{sql.NullFloat64{Float64: 1.23, Valid: true}},
	}
}

type MeasurementsService interface {
	// Get a measurement
	Get(id int64) (*Measurement, error)

	// List all measurements
	List(opt *MeasurementListOptions) ([]*Measurement, error)

	// Create a measurement
	Create(measurement *Measurement) (bool, error)

	// Update an existing measurement
	Update(id int64, MeasurementType *Measurement) (bool, error)

	// Delete a measurement
	Delete(id int64) (deleted bool, err error)
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

	var measurement *MeasurementJSON
	_, err = s.client.Do(req, &measurement)
	if err != nil {
		return nil, err
	}

	return measurement.Measurement, nil
}

func (s *measurementsService) Create(measurement *Measurement) (bool, error) {
	url, err := s.client.url(router.CreateMeasurement, nil, nil)
	if err != nil {
		return false, err
	}

	req, err := s.client.NewRequest("POST", url.String(), MeasurementJSON{Measurement: measurement})
	if err != nil {
		return false, err
	}

	resp, err := s.client.Do(req, &measurement)
	if err != nil {
		return false, err
	}

	return resp.StatusCode == http.StatusCreated, nil
}

type MeasurementListOptions struct {
	ListOptions
	Genus string
}

func (s *measurementsService) List(opt *MeasurementListOptions) ([]*Measurement, error) {
	url, err := s.client.url(router.Measurements, nil, opt)
	if err != nil {
		return nil, err
	}

	req, err := s.client.NewRequest("GET", url.String(), nil)
	if err != nil {
		return nil, err
	}

	var measurements *MeasurementsJSON
	_, err = s.client.Do(req, &measurements)
	if err != nil {
		return nil, err
	}

	return measurements.Measurements, nil
}

func (s *measurementsService) Update(id int64, measurement *Measurement) (bool, error) {
	strId := strconv.FormatInt(id, 10)

	url, err := s.client.url(router.UpdateMeasurement, map[string]string{"Id": strId}, nil)
	if err != nil {
		return false, err
	}

	req, err := s.client.NewRequest("PUT", url.String(), MeasurementJSON{Measurement: measurement})
	if err != nil {
		return false, err
	}

	resp, err := s.client.Do(req, &measurement)
	if err != nil {
		return false, err
	}

	return resp.StatusCode == http.StatusOK, nil
}

func (s *measurementsService) Delete(id int64) (bool, error) {
	strId := strconv.FormatInt(id, 10)

	url, err := s.client.url(router.DeleteMeasurement, map[string]string{"Id": strId}, nil)
	if err != nil {
		return false, err
	}

	req, err := s.client.NewRequest("DELETE", url.String(), nil)
	if err != nil {
		return false, err
	}

	var measurement *Measurement
	resp, err := s.client.Do(req, &measurement)
	if err != nil {
		return false, err
	}

	return resp.StatusCode == http.StatusOK, nil
}

type MockMeasurementsService struct {
	Get_    func(id int64) (*Measurement, error)
	List_   func(opt *MeasurementListOptions) ([]*Measurement, error)
	Create_ func(measurement *Measurement) (bool, error)
	Update_ func(id int64, measurement *Measurement) (bool, error)
	Delete_ func(id int64) (bool, error)
}

var _ MeasurementsService = &MockMeasurementsService{}

func (s *MockMeasurementsService) Get(id int64) (*Measurement, error) {
	if s.Get_ == nil {
		return nil, nil
	}
	return s.Get_(id)
}

func (s *MockMeasurementsService) Create(measurement *Measurement) (bool, error) {
	if s.Create_ == nil {
		return false, nil
	}
	return s.Create_(measurement)
}

func (s *MockMeasurementsService) List(opt *MeasurementListOptions) ([]*Measurement, error) {
	if s.List_ == nil {
		return nil, nil
	}
	return s.List_(opt)
}

func (s *MockMeasurementsService) Update(id int64, measurement *Measurement) (bool, error) {
	if s.Update_ == nil {
		return false, nil
	}
	return s.Update_(id, measurement)
}

func (s *MockMeasurementsService) Delete(id int64) (bool, error) {
	if s.Delete_ == nil {
		return false, nil
	}
	return s.Delete_(id)
}
