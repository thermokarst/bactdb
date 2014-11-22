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

func (s *observationsStore) Update(id int64, observation *models.Observation) (bool, error) {
	_, err := s.Get(id)
	if err != nil {
		return false, err
	}

	if id != observation.Id {
		return false, models.ErrObservationNotFound
	}

	observation.UpdatedAt = time.Now()
	changed, err := s.dbh.Update(observation)
	if err != nil {
		return false, err
	}

	if changed == 0 {
		return false, ErrNoRowsUpdated
	}

	return true, nil
}

func (s *observationsStore) Delete(id int64) (bool, error) {
	observation, err := s.Get(id)
	if err != nil {
		return false, err
	}

	deleted, err := s.dbh.Delete(observation)
	if err != nil {
		return false, err
	}
	if deleted == 0 {
		return false, ErrNoRowsDeleted
	}
	return true, nil
}
