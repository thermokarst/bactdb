package datastore

import (
	"time"

	"github.com/thermokarst/bactdb/models"
)

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

func (s *observationTypesStore) Create(observation_type *models.ObservationType) (bool, error) {
	currentTime := time.Now()
	observation_type.CreatedAt = currentTime
	observation_type.UpdatedAt = currentTime
	if err := s.dbh.Insert(observation_type); err != nil {
		return false, err
	}
	return true, nil
}

func (s *observationTypesStore) List(opt *models.ObservationTypeListOptions) ([]*models.ObservationType, error) {
	if opt == nil {
		opt = &models.ObservationTypeListOptions{}
	}
	var observation_types []*models.ObservationType
	err := s.dbh.Select(&observation_types, `SELECT * FROM observation_types LIMIT $1 OFFSET $2;`, opt.PerPageOrDefault(), opt.Offset())
	if err != nil {
		return nil, err
	}
	return observation_types, nil
}