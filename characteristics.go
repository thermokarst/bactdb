package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"strings"
	"time"
)

func init() {
	DB.AddTableWithName(CharacteristicBase{}, "characteristics").SetKeys(true, "Id")
}

type CharacteristicService struct{}

// A Characteristic is a lookup type
type CharacteristicBase struct {
	Id                   int64     `json:"id,omitempty"`
	CharacteristicName   string    `db:"characteristic_name" json:"characteristicName"`
	CharacteristicTypeId int64     `db:"characteristic_type_id" json:"characteristicType"`
	SortOrder            NullInt64 `db:"sort_order" json:"sortOrder"`
	CreatedAt            time.Time `db:"created_at" json:"createdAt"`
	UpdatedAt            time.Time `db:"updated_at" json:"updatedAt"`
	DeletedAt            NullTime  `db:"deleted_at" json:"deletedAt"`
	CreatedBy            int64     `db:"created_by" json:"createdBy"`
	UpdatedBy            int64     `db:"updated_by" json:"updatedBy"`
	DeletedBy            NullInt64 `db:"deleted_by" json:"deletedBy"`
}

type Characteristic struct {
	*CharacteristicBase
	Measurements NullSliceInt64 `db:"measurements" json:"measurements"`
	Strains      NullSliceInt64 `db:"strains" json:"strains"`
}

type Characteristics []*Characteristic

type CharacteristicJSON struct {
	Characteristic *Characteristic `json:"characteristic"`
}

type CharacteristicsJSON struct {
	Characteristics *Characteristics `json:"characteristics"`
}

func (c *Characteristic) marshal() ([]byte, error) {
	return json.Marshal(&CharacteristicJSON{Characteristic: c})
}

func (c *Characteristics) marshal() ([]byte, error) {
	return json.Marshal(&CharacteristicsJSON{Characteristics: c})
}

func (c CharacteristicService) list(val *url.Values) (entity, error) {
	if val == nil {
		return nil, errors.New("must provide options")
	}
	var opt ListOptions
	if err := schemaDecoder.Decode(&opt, *val); err != nil {
		return nil, err
	}

	var vals []interface{}
	sql := `SELECT c.*, array_agg(m.id) AS measurements, array_agg(st.id) AS strains
			FROM characteristics c
			LEFT OUTER JOIN measurements m ON m.characteristic_id=c.id
			LEFT OUTER JOIN strains st ON st.id=m.strain_id`

	if len(opt.Ids) != 0 {
		var conds []string

		c := "c.id IN ("
		for i, id := range opt.Ids {
			c = c + fmt.Sprintf("$%v,", i+1) // start param index at 1
			vals = append(vals, id)
		}
		c = c[:len(c)-1] + ")"
		conds = append(conds, c)
		sql += " WHERE (" + strings.Join(conds, ") AND (") + ")"
	}

	sql += " GROUP BY c.id;"

	characteristics := make(Characteristics, 0)
	err := DBH.Select(&characteristics, sql, vals...)
	if err != nil {
		return nil, err
	}
	return &characteristics, nil
}

func (c CharacteristicService) get(id int64, dummy string) (entity, error) {
	var characteristic Characteristic
	sql := `SELECT c.*, array_agg(m.id) AS measurements, array_agg(st.id) AS strains
			FROM characteristics c
			LEFT OUTER JOIN measurements m ON m.characteristic_id=c.id
			LEFT OUTER JOIN strains st ON st.id=m.strain_id
			WHERE c.id=$1
			GROUP BY c.id;`
	if err := DBH.SelectOne(&characteristic, sql, id); err != nil {
		return nil, err
	}
	return &characteristic, nil
}
