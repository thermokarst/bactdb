package datastore

import (
	"time"

	"github.com/thermokarst/bactdb/models"
)

func init() {
	DB.AddTableWithName(models.TextMeasurementType{}, "text_measurement_types").SetKeys(true, "Id")
}

type textMeasurementTypesStore struct {
	*Datastore
}

func (s *textMeasurementTypesStore) Get(id int64) (*models.TextMeasurementType, error) {
	var text_measurement_type []*models.TextMeasurementType
	if err := s.dbh.Select(&text_measurement_type, `SELECT * FROM text_measurement_types WHERE id=$1;`, id); err != nil {
		return nil, err
	}
	if len(text_measurement_type) == 0 {
		return nil, models.ErrTextMeasurementTypeNotFound
	}
	return text_measurement_type[0], nil
}

func (s *textMeasurementTypesStore) Create(text_measurement_type *models.TextMeasurementType) (bool, error) {
	currentTime := time.Now()
	text_measurement_type.CreatedAt = currentTime
	text_measurement_type.UpdatedAt = currentTime
	if err := s.dbh.Insert(text_measurement_type); err != nil {
		return false, err
	}
	return true, nil
}

func (s *textMeasurementTypesStore) List(opt *models.TextMeasurementTypeListOptions) ([]*models.TextMeasurementType, error) {
	if opt == nil {
		opt = &models.TextMeasurementTypeListOptions{}
	}
	var text_measurement_types []*models.TextMeasurementType
	err := s.dbh.Select(&text_measurement_types, `SELECT * FROM text_measurement_types LIMIT $1 OFFSET $2;`, opt.PerPageOrDefault(), opt.Offset())
	if err != nil {
		return nil, err
	}
	return text_measurement_types, nil
}

func (s *textMeasurementTypesStore) Update(id int64, text_measurement_type *models.TextMeasurementType) (bool, error) {
	_, err := s.Get(id)
	if err != nil {
		return false, err
	}

	if id != text_measurement_type.Id {
		return false, models.ErrObservationNotFound
	}

	text_measurement_type.UpdatedAt = time.Now()
	changed, err := s.dbh.Update(text_measurement_type)
	if err != nil {
		return false, err
	}

	if changed == 0 {
		return false, ErrNoRowsUpdated
	}

	return true, nil
}

func (s *textMeasurementTypesStore) Delete(id int64) (bool, error) {
	text_measurement_type, err := s.Get(id)
	if err != nil {
		return false, err
	}

	deleted, err := s.dbh.Delete(text_measurement_type)
	if err != nil {
		return false, err
	}
	if deleted == 0 {
		return false, ErrNoRowsDeleted
	}
	return true, nil
}
