package datastore

import (
	"strings"
	"time"

	"github.com/thermokarst/bactdb/models"
)

func init() {
	DB.AddTableWithName(models.Species{}, "species").SetKeys(true, "Id")
}

type speciesStore struct {
	*Datastore
}

func (s *speciesStore) Get(id int64) (*models.Species, error) {
	var species models.Species
	if err := s.dbh.SelectOne(&species, `SELECT * FROM species WHERE id=$1;`, id); err != nil {
		return nil, err
	}
	if &species == nil {
		return nil, models.ErrSpeciesNotFound
	}
	return &species, nil
}

func (s *speciesStore) Create(species *models.Species) (bool, error) {
	currentTime := time.Now()
	species.CreatedAt = currentTime
	species.UpdatedAt = currentTime
	if err := s.dbh.Insert(species); err != nil {
		return false, err
	}
	return true, nil
}

func (s *speciesStore) List(opt *models.SpeciesListOptions) ([]*models.Species, error) {
	if opt == nil {
		opt = &models.SpeciesListOptions{}
	}

	sql := `SELECT * FROM species`

	var conds []string
	var vals []interface{}

	if opt.Genus != "" {
		conds = append(conds, "genus_id = (SELECT id FROM genera WHERE lower(genus_name) = $1)")
		vals = append(vals, opt.Genus)
	}

	if len(conds) > 0 {
		sql += " WHERE (" + strings.Join(conds, ") AND (") + ")"
	}

	sql += ";"

	var species []*models.Species
	err := s.dbh.Select(&species, sql, vals...)
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

	species.UpdatedAt = time.Now()
	changed, err := s.dbh.Update(species)
	if err != nil {
		return false, err
	}

	if changed == 0 {
		return false, ErrNoRowsUpdated
	}

	return true, nil
}

func (s *speciesStore) Delete(id int64) (bool, error) {
	species, err := s.Get(id)
	if err != nil {
		return false, err
	}

	deleted, err := s.dbh.Delete(species)
	if err != nil {
		return false, err
	}
	if deleted == 0 {
		return false, ErrNoRowsDeleted
	}
	return true, nil
}
