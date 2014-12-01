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

func (s *unitTypesStore) Update(id int64, unit_type *models.UnitType) (bool, error) {
	_, err := s.Get(id)
	if err != nil {
		return false, err
	}

	if id != unit_type.Id {
		return false, models.ErrUnitTypeNotFound
	}

	unit_type.UpdatedAt = time.Now()
	changed, err := s.dbh.Update(unit_type)
	if err != nil {
		return false, err
	}

	if changed == 0 {
		return false, ErrNoRowsUpdated
	}

	return true, nil
}

func (s *unitTypesStore) Delete(id int64) (bool, error) {
	unit_type, err := s.Get(id)
	if err != nil {
		return false, err
	}

	deleted, err := s.dbh.Delete(unit_type)
	if err != nil {
		return false, err
	}
	if deleted == 0 {
		return false, ErrNoRowsDeleted
	}
	return true, nil
}