package datastore

import "github.com/thermokarst/bactdb/models"

func init() {
	DB.AddTableWithName(models.Species{}, "species").SetKeys(true, "Id")
}

type speciesStore struct {
	*Datastore
}

func (s *speciesStore) Get(id int64) (*models.Species, error) {
	var species []*models.Species
	if err := s.dbh.Select(&species, `SELECT * FROM species WHERE id=$1;`, id); err != nil {
		return nil, err
	}
	if len(species) == 0 {
		return nil, models.ErrSpeciesNotFound
	}
	return species[0], nil
}

func (s *speciesStore) Create(species *models.Species) (bool, error) {
	if err := s.dbh.Insert(species); err != nil {
		return false, err
	}
	return true, nil
}

func (s *speciesStore) List(opt *models.SpeciesListOptions) ([]*models.Species, error) {
	if opt == nil {
		opt = &models.SpeciesListOptions{}
	}
	var species []*models.Species
	err := s.dbh.Select(&species, `SELECT * FROM species LIMIT $1 OFFSET $2;`, opt.PerPageOrDefault(), opt.Offset())
	if err != nil {
		return nil, err
	}
	return species, nil
}

func (s *speciesStore) Update(id int64, species *models.Species) (bool, error) {
	_, err := s.Get(id)
	if err != nil {
		return false, err
	}

	if id != species.Id {
		return false, models.ErrSpeciesNotFound
	}

	changed, err := s.dbh.Update(species)
	if err != nil {
		return false, err
	}

	if changed == 0 {
		return false, ErrNoRowsUpdated
	}

	return true, nil
}
