package datastore

import "github.com/thermokarst/bactdb/models"

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
