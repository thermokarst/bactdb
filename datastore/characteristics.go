package datastore

import (
	"time"

	"github.com/thermokarst/bactdb/models"
)

func init() {
	DB.AddTableWithName(models.Characteristic{}, "characteristics").SetKeys(true, "Id")
}

type characteristicsStore struct {
	*Datastore
}

func (s *characteristicsStore) Get(id int64) (*models.Characteristic, error) {
	var characteristic models.Characteristic
	if err := s.dbh.SelectOne(&characteristic, `SELECT * FROM characteristics WHERE id=$1;`, id); err != nil {
		return nil, err
	}
	if &characteristic == nil {
		return nil, models.ErrCharacteristicNotFound
	}
	return &characteristic, nil
}

func (s *characteristicsStore) Create(characteristic *models.Characteristic) (bool, error) {
	currentTime := time.Now()
	characteristic.CreatedAt = currentTime
	characteristic.UpdatedAt = currentTime
	if err := s.dbh.Insert(characteristic); err != nil {
		return false, err
	}
	return true, nil
}

func (s *characteristicsStore) List(opt *models.CharacteristicListOptions) ([]*models.Characteristic, error) {
	if opt == nil {
		opt = &models.CharacteristicListOptions{}
	}
	var characteristics []*models.Characteristic
	err := s.dbh.Select(&characteristics, `SELECT * FROM characteristics LIMIT $1 OFFSET $2;`, opt.PerPageOrDefault(), opt.Offset())
	if err != nil {
		return nil, err
	}
	return characteristics, nil
}

func (s *characteristicsStore) Update(id int64, characteristic *models.Characteristic) (bool, error) {
	_, err := s.Get(id)
	if err != nil {
		return false, err
	}

	if id != characteristic.Id {
		return false, models.ErrCharacteristicNotFound
	}

	characteristic.UpdatedAt = time.Now()
	changed, err := s.dbh.Update(characteristic)
	if err != nil {
		return false, err
	}

	if changed == 0 {
		return false, ErrNoRowsUpdated
	}

	return true, nil
}

func (s *characteristicsStore) Delete(id int64) (bool, error) {
	characteristic, err := s.Get(id)
	if err != nil {
		return false, err
	}

	deleted, err := s.dbh.Delete(characteristic)
	if err != nil {
		return false, err
	}
	if deleted == 0 {
		return false, ErrNoRowsDeleted
	}
	return true, nil
}
