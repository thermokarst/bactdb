package datastore

import "github.com/thermokarst/bactdb/models"

func init() {
	DB.AddTableWithName(models.Strain{}, "strains").SetKeys(true, "Id")
}

type strainsStore struct {
	*Datastore
}

func (s *strainsStore) Get(id int64) (*models.Strain, error) {
	var strain []*models.Strain
	if err := s.dbh.Select(&strain, `SELECT * FROM strains WHERE id=$1;`, id); err != nil {
		return nil, err
	}
	if len(strain) == 0 {
		return nil, models.ErrStrainNotFound
	}
	return strain[0], nil
}

func (s *strainsStore) Create(strain *models.Strain) (bool, error) {
	if err := s.dbh.Insert(strain); err != nil {
		return false, err
	}
	return true, nil
}
