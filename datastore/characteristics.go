package datastore

import (
	"time"

	"github.com/thermokarst/bactdb/models"
)

func init() {
	DB.AddTableWithName(models.CharacteristicBase{}, "characteristics").SetKeys(true, "Id")
}

type characteristicsStore struct {
	*Datastore
}

func (s *characteristicsStore) Get(id int64) (*models.Characteristic, error) {
	var characteristic models.Characteristic
	err := s.dbh.SelectOne(&characteristic, `SELECT c.*, array_agg(m.id) AS measurements FROM characteristics c LEFT OUTER JOIN measurements m ON m.characteristic_id=c.id WHERE c.id=$1 GROUP BY c.id;`, id)
	if err != nil {
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
	base := characteristic.CharacteristicBase
	if err := s.dbh.Insert(base); err != nil {
		return false, err
	}
	characteristic.Id = base.Id
	return true, nil
}

func (s *characteristicsStore) List(opt *models.CharacteristicListOptions) ([]*models.Characteristic, error) {
	if opt == nil {
		opt = &models.CharacteristicListOptions{}
	}
	var characteristics []*models.Characteristic
	err := s.dbh.Select(&characteristics, `SELECT c.*, array_agg(m.id) AS measurements FROM characteristics c LEFT OUTER JOIN measurements m ON m.characteristic_id=c.id GROUP BY c.id;`)
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
	changed, err := s.dbh.Update(characteristic.CharacteristicBase)
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

	deleted, err := s.dbh.Delete(characteristic.CharacteristicBase)
	if err != nil {
		return false, err
	}
	if deleted == 0 {
		return false, ErrNoRowsDeleted
	}
	return true, nil
}
