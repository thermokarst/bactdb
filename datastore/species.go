package datastore

import (
	"strings"
	"time"

	"github.com/thermokarst/bactdb/models"
)

func init() {
	DB.AddTableWithName(models.SpeciesBase{}, "species").SetKeys(true, "Id")
}

type speciesStore struct {
	*Datastore
}

func (s *speciesStore) Get(id int64) (*models.Species, error) {
	var species models.Species
	err := s.dbh.SelectOne(&species, `SELECT sp.*, array_agg(st.id) AS strains FROM species sp LEFT OUTER JOIN strains st ON st.species_id=sp.id WHERE sp.id=$1 GROUP BY sp.id;`, id)
	if err != nil {
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
	base := species.SpeciesBase
	if err := s.dbh.Insert(base); err != nil {
		return false, err
	}
	species.Id = base.Id
	return true, nil
}

func (s *speciesStore) List(opt *models.SpeciesListOptions) ([]*models.Species, error) {
	if opt == nil {
		opt = &models.SpeciesListOptions{}
	}

	sql := `SELECT sp.*, array_agg(st.id) AS strains FROM species sp LEFT OUTER JOIN strains st ON st.species_id=sp.id`

	var conds []string
	var vals []interface{}

	if opt.Genus != "" {
		conds = append(conds, "sp.genus_id = (SELECT id FROM genera WHERE lower(genus_name) = $1)")
		vals = append(vals, opt.Genus)
	}

	if len(conds) > 0 {
		sql += " WHERE (" + strings.Join(conds, ") AND (") + ")"
	}

	sql += " GROUP BY sp.id;"

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

	changed, err := s.dbh.Update(species.SpeciesBase)
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

	deleted, err := s.dbh.Delete(species.SpeciesBase)
	if err != nil {
		return false, err
	}
	if deleted == 0 {
		return false, ErrNoRowsDeleted
	}
	return true, nil
}
