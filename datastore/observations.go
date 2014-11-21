package datastore

import (
	"time"

	"github.com/thermokarst/bactdb/models"
)

func init() {
	DB.AddTableWithName(models.Observation{}, "observations").SetKeys(true, "Id")
}

type observationsStore struct {
	*Datastore
}

func (s *observationsStore) Get(id int64) (*models.Observation, error) {
	var observation []*models.Observation
	if err := s.dbh.Select(&observation, `SELECT * FROM observations WHERE id=$1;`, id); err != nil {
		return nil, err
	}
	if len(observation) == 0 {
		return nil, models.ErrObservationNotFound
	}
	return observation[0], nil
}

func (s *observationsStore) Create(observation *models.Observation) (bool, error) {
	currentTime := time.Now()
	observation.CreatedAt = currentTime
	observation.UpdatedAt = currentTime
	if err := s.dbh.Insert(observation); err != nil {
		return false, err
	}
	return true, nil
}

func (s *observationsStore) List(opt *models.ObservationListOptions) ([]*models.Observation, error) {
	if opt == nil {
		opt = &models.ObservationListOptions{}
	}
	var observations []*models.Observation
	err := s.dbh.Select(&observations, `SELECT * FROM observations LIMIT $1 OFFSET $2;`, opt.PerPageOrDefault(), opt.Offset())
	if err != nil {
		return nil, err
	}
	return observations, nil
}
