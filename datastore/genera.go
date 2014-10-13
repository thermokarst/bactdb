package datastore

import (
	"strings"

	"github.com/thermokarst/bactdb/models"
)

func init() {
	DB.AddTableWithName(models.Genus{}, "genera").SetKeys(true, "Id")
}

type generaStore struct {
	*Datastore
}

func (s *generaStore) Get(id int64) (*models.Genus, error) {
	var genus []*models.Genus
	if err := s.dbh.Select(&genus, `SELECT * FROM genera WHERE id=$1;`, id); err != nil {
		return nil, err
	}
	if len(genus) == 0 {
		return nil, models.ErrGenusNotFound
	}
	return genus[0], nil
}

func (s *generaStore) Create(genus *models.Genus) (bool, error) {
	if err := s.dbh.Insert(genus); err != nil {
		if strings.Contains(err.Error(), `violates unique constraint "genus_idx"`) {
			return false, err
		}
	}
	return true, nil
}

func (s *generaStore) List(opt *models.GenusListOptions) ([]*models.Genus, error) {
	if opt == nil {
		opt = &models.GenusListOptions{}
	}
	var genera []*models.Genus
	err := s.dbh.Select(&genera, `SELECT * FROM genera LIMIT $1 OFFSET $2;`, opt.PerPageOrDefault(), opt.Offset())
	if err != nil {
		return nil, err
	}
	return genera, nil
}
