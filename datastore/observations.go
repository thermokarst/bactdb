package datastore

import "github.com/thermokarst/bactdb/models"

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
