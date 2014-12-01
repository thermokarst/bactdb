package datastore

import (
	"time"

	"github.com/thermokarst/bactdb/models"
)

func init() {
	DB.AddTableWithName(models.UnitType{}, "unit_types").SetKeys(true, "Id")
}

type unitTypesStore struct {
	*Datastore
}

func (s *unitTypesStore) Get(id int64) (*models.UnitType, error) {
	var unit_type []*models.UnitType
	if err := s.dbh.Select(&unit_type, `SELECT * FROM unit_types WHERE id=$1;`, id); err != nil {
		return nil, err
	}
	if len(unit_type) == 0 {
		return nil, models.ErrUnitTypeNotFound
	}
	return unit_type[0], nil
}

func (s *unitTypesStore) Create(unit_type *models.UnitType) (bool, error) {
	currentTime := time.Now()
	unit_type.CreatedAt = currentTime
	unit_type.UpdatedAt = currentTime
	if err := s.dbh.Insert(unit_type); err != nil {
		return false, err
	}
	return true, nil
}

func (s *unitTypesStore) List(opt *models.UnitTypeListOptions) ([]*models.UnitType, error) {
	if opt == nil {
		opt = &models.UnitTypeListOptions{}
	}
	var unit_types []*models.UnitType
	err := s.dbh.Select(&unit_types, `SELECT * FROM unit_types LIMIT $1 OFFSET $2;`, opt.PerPageOrDefault(), opt.Offset())
	if err != nil {
		return nil, err
	}
	return unit_types, nil
}
