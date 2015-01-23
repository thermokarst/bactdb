package datastore

import (
	"fmt"
	"strings"
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
	var measurement models.Measurement
	if err := s.dbh.SelectOne(&measurement, `SELECT * FROM measurements WHERE id=$1;`, id); err != nil {
		return nil, err
	}
	if &measurement == nil {
		return nil, models.ErrMeasurementNotFound
	}
	return &measurement, nil
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

	sql := `SELECT * FROM measurements`

	var conds []string
	var vals []interface{}

	if opt.Genus != "" {
		conds = append(conds, `strain_id IN (SELECT st.id FROM strains st
		INNER JOIN species sp ON sp.id = st.species_id
		INNER JOIN genera g ON g.id = sp.genus_id
		WHERE lower(g.genus_name) = $1)`)
		vals = append(vals, opt.Genus)
	}

	if len(conds) > 0 {
		sql += " WHERE (" + strings.Join(conds, ") AND (") + ")"
	}

	sql += fmt.Sprintf(" LIMIT $%v OFFSET $%v;", len(conds)+1, len(conds)+2)
	vals = append(vals, opt.PerPageOrDefault(), opt.Offset())

	var measurements []*models.Measurement
	err := s.dbh.Select(&measurements, sql, vals...)
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

func (s *measurementsStore) Delete(id int64) (bool, error) {
	measurement, err := s.Get(id)
	if err != nil {
		return false, err
	}

	deleted, err := s.dbh.Delete(measurement)
	if err != nil {
		return false, err
	}
	if deleted == 0 {
		return false, ErrNoRowsDeleted
	}
	return true, nil
}
