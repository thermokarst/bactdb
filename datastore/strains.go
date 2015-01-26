package datastore

import (
	"strings"
	"time"

	"github.com/thermokarst/bactdb/models"
)

func init() {
	DB.AddTableWithName(models.StrainBase{}, "strains").SetKeys(true, "Id")
}

type strainsStore struct {
	*Datastore
}

func (s *strainsStore) Get(id int64) (*models.Strain, error) {
	var strain models.Strain
	err := s.dbh.SelectOne(&strain, `SELECT s.*, array_agg(m.id) AS measurements FROM strains s LEFT OUTER JOIN measurements m ON m.strain_id=s.id WHERE s.id=$1 GROUP BY s.id;`, id)
	if err != nil {
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
	base := strain.StrainBase
	if err := s.dbh.Insert(base); err != nil {
		return false, err
	}
	strain.Id = base.Id
	return true, nil
}

func (s *strainsStore) List(opt *models.StrainListOptions) ([]*models.Strain, error) {
	if opt == nil {
		opt = &models.StrainListOptions{}
	}

	sql := `SELECT s.*, array_agg(m.id) AS measurements FROM strains s LEFT OUTER JOIN measurements m ON m.strain_id=s.id`

	var conds []string
	var vals []interface{}

	if opt.Genus != "" {
		conds = append(conds, `s.species_id IN (SELECT s.id FROM species s
		INNER JOIN genera g ON g.id = s.genus_id WHERE lower(g.genus_name) = $1)`)
		vals = append(vals, opt.Genus)
	}

	if len(conds) > 0 {
		sql += " WHERE (" + strings.Join(conds, ") AND (") + ")"
	}

	sql += " GROUP BY s.id;"

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

	changed, err := s.dbh.Update(strain.StrainBase)
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

	deleted, err := s.dbh.Delete(strain.StrainBase)
	if err != nil {
		return false, err
	}
	if deleted == 0 {
		return false, ErrNoRowsDeleted
	}
	return true, nil
}
