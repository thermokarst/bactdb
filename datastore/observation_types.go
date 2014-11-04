package datastore

import "github.com/thermokarst/bactdb/models"

func init() {
	DB.AddTableWithName(models.ObservationType{}, "observation_types").SetKeys(true, "Id")
}

type observationTypesStore struct {
	*Datastore
}

func (s *observationTypesStore) Get(id int64) (*models.ObservationType, error) {
	var observation_type []*models.ObservationType
	if err := s.dbh.Select(&observation_type, `SELECT * FROM observation_types WHERE id=$1;`, id); err != nil {
		return nil, err
	}
	if len(observation_type) == 0 {
		return nil, models.ErrObservationTypeNotFound
	}
	return observation_type[0], nil
}
