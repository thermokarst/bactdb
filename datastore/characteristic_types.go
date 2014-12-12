package datastore

import (
	"time"

	"github.com/thermokarst/bactdb/models"
)

func init() {
	DB.AddTableWithName(models.CharacteristicType{}, "characteristic_types").SetKeys(true, "Id")
}

type characteristicTypesStore struct {
	*Datastore
}

func (s *characteristicTypesStore) Get(id int64) (*models.CharacteristicType, error) {
	var characteristic_type []*models.CharacteristicType
	if err := s.dbh.Select(&characteristic_type, `SELECT * FROM characteristic_types WHERE id=$1;`, id); err != nil {
		return nil, err
	}
	if len(characteristic_type) == 0 {
		return nil, models.ErrCharacteristicTypeNotFound
	}
	return characteristic_type[0], nil
}

func (s *characteristicTypesStore) Create(characteristic_type *models.CharacteristicType) (bool, error) {
	currentTime := time.Now()
	characteristic_type.CreatedAt = currentTime
	characteristic_type.UpdatedAt = currentTime
	if err := s.dbh.Insert(characteristic_type); err != nil {
		return false, err
	}
	return true, nil
}

func (s *characteristicTypesStore) List(opt *models.CharacteristicTypeListOptions) ([]*models.CharacteristicType, error) {
	if opt == nil {
		opt = &models.CharacteristicTypeListOptions{}
	}
	var characteristic_types []*models.CharacteristicType
	err := s.dbh.Select(&characteristic_types, `SELECT * FROM characteristic_types LIMIT $1 OFFSET $2;`, opt.PerPageOrDefault(), opt.Offset())
	if err != nil {
		return nil, err
	}
	return characteristic_types, nil
}

func (s *characteristicTypesStore) Update(id int64, characteristic_type *models.CharacteristicType) (bool, error) {
	_, err := s.Get(id)
	if err != nil {
		return false, err
	}

	if id != characteristic_type.Id {
		return false, models.ErrCharacteristicTypeNotFound
	}

	characteristic_type.UpdatedAt = time.Now()
	changed, err := s.dbh.Update(characteristic_type)
	if err != nil {
		return false, err
	}

	if changed == 0 {
		return false, ErrNoRowsUpdated
	}

	return true, nil
}

func (s *characteristicTypesStore) Delete(id int64) (bool, error) {
	characteristic_type, err := s.Get(id)
	if err != nil {
		return false, err
	}

	deleted, err := s.dbh.Delete(characteristic_type)
	if err != nil {
		return false, err
	}
	if deleted == 0 {
		return false, ErrNoRowsDeleted
	}
	return true, nil
}
