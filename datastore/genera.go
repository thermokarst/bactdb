package datastore

import (
	"strings"
	"time"

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
	currentTime := time.Now()
	genus.CreatedAt = currentTime
	genus.UpdatedAt = currentTime
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

func (s *generaStore) Update(id int64, genus *models.Genus) (bool, error) {
	_, err := s.Get(id)
	if err != nil {
		return false, err
	}

	if id != genus.Id {
		return false, models.ErrGenusNotFound
	}

	genus.UpdatedAt = time.Now()
	changed, err := s.dbh.Update(genus)
	if err != nil {
		return false, err
	}

	if changed == 0 {
		return false, ErrNoRowsUpdated
	}

	return true, nil
}

func (s *generaStore) Delete(id int64) (bool, error) {
	genus, err := s.Get(id)
	if err != nil {
		return false, err
	}

	deleted, err := s.dbh.Delete(genus)
	if err != nil {
		return false, err
	}
	if deleted == 0 {
		return false, ErrNoRowsDeleted
	}
	return true, nil
}
