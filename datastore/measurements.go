package datastore

import (
	"time"

	"github.com/thermokarst/bactdb/models"
)

func init() {
	DB.AddTableWithName(models.Measurement{}, "measurements").SetKeys(true, "Id")
}

type measurementsStore struct {
	*Datastore
}

func (s *measurementsStore) Get(id int64) (*models.Measurement, error) {
	var measurement []*models.Measurement
	if err := s.dbh.Select(&measurement, `SELECT * FROM measurements WHERE id=$1;`, id); err != nil {
		return nil, err
	}
	if len(measurement) == 0 {
		return nil, models.ErrMeasurementNotFound
	}
	return measurement[0], nil
}

func (s *measurementsStore) Create(measurement *models.Measurement) (bool, error) {
	currentTime := time.Now()
	measurement.CreatedAt = currentTime
	measurement.UpdatedAt = currentTime
	if err := s.dbh.Insert(measurement); err != nil {
		return false, err
	}
	return true, nil
}

func (s *measurementsStore) List(opt *models.MeasurementListOptions) ([]*models.Measurement, error) {
	if opt == nil {
		opt = &models.MeasurementListOptions{}
	}
	var measurements []*models.Measurement
	err := s.dbh.Select(&measurements, `SELECT * FROM measurements LIMIT $1 OFFSET $2;`, opt.PerPageOrDefault(), opt.Offset())
	if err != nil {
		return nil, err
	}
	return measurements, nil
}

func (s *measurementsStore) Update(id int64, measurement *models.Measurement) (bool, error) {
	_, err := s.Get(id)
	if err != nil {
		return false, err
	}

	if id != measurement.Id {
		return false, models.ErrMeasurementNotFound
	}

	measurement.UpdatedAt = time.Now()
	changed, err := s.dbh.Update(measurement)
	if err != nil {
		return false, err
	}

	if changed == 0 {
		return false, ErrNoRowsUpdated
	}

	return true, nil
}
