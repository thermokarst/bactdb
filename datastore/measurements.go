package datastore

import "github.com/thermokarst/bactdb/models"

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
