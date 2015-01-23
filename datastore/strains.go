package datastore

import (
	"fmt"
	"strings"
	"time"

	"github.com/thermokarst/bactdb/models"
)

func init() {
	DB.AddTableWithName(models.Strain{}, "strains").SetKeys(true, "Id")
}

type strainsStore struct {
	*Datastore
}

func (s *strainsStore) Get(id int64) (*models.Strain, error) {
	var strain models.Strain
	if err := s.dbh.SelectOne(&strain, `SELECT * FROM strains WHERE id=$1;`, id); err != nil {
		return nil, err
	}
	if &strain == nil {
		return nil, models.ErrStrainNotFound
	}
	return &strain, nil
}

func (s *strainsStore) Create(strain *models.Strain) (bool, error) {
	currentTime := time.Now()
	strain.CreatedAt = currentTime
	strain.UpdatedAt = currentTime
	if err := s.dbh.Insert(strain); err != nil {
		return false, err
	}
	return true, nil
}

func (s *strainsStore) List(opt *models.StrainListOptions) ([]*models.Strain, error) {
	if opt == nil {
		opt = &models.StrainListOptions{}
	}

	sql := `SELECT * FROM strains`

	var conds []string
	var vals []interface{}

	if opt.Genus != "" {
		conds = append(conds, `species_id IN (SELECT s.id FROM species s
		INNER JOIN genera g ON g.id = s.genus_id WHERE lower(g.genus_name) = $1)`)
		vals = append(vals, opt.Genus)
	}

	if len(conds) > 0 {
		sql += " WHERE (" + strings.Join(conds, ") AND (") + ")"
	}

	sql += fmt.Sprintf(" LIMIT $%v OFFSET $%v;", len(conds)+1, len(conds)+2)
	vals = append(vals, opt.PerPageOrDefault(), opt.Offset())

	var strains []*models.Strain
	err := s.dbh.Select(&strains, sql, vals...)
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

	strain.UpdatedAt = time.Now()
	changed, err := s.dbh.Update(strain)
	if err != nil {
		return false, err
	}

	if changed == 0 {
		return false, ErrNoRowsUpdated
	}

	return true, nil
}

func (s *strainsStore) Delete(id int64) (bool, error) {
	strain, err := s.Get(id)
	if err != nil {
		return false, err
	}

	deleted, err := s.dbh.Delete(strain)
	if err != nil {
		return false, err
	}
	if deleted == 0 {
		return false, ErrNoRowsDeleted
	}
	return true, nil
}
