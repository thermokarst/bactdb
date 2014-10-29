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

func (s *strainsStore) List(opt *models.StrainListOptions) ([]*models.Strain, error) {
	if opt == nil {
		opt = &models.StrainListOptions{}
	}
	var strains []*models.Strain
	err := s.dbh.Select(&strains, `SELECT * FROM strains LIMIT $1 OFFSET $2;`, opt.PerPageOrDefault(), opt.Offset())
	if err != nil {
		return nil, err
	}
	return strains, nil
}

func (s *strainsStore) Update(id int64, strain *models.Strain) (bool, error) {
	_, err := s.Get(id)
	if err != nil {
		return false, err
	}

	if id != strain.Id {
		return false, models.ErrStrainNotFound
	}

	changed, err := s.dbh.Update(strain)
	if err != nil {
		return false, err
	}

	if changed == 0 {
		return false, ErrNoRowsUpdated
	}

	return true, nil
}
