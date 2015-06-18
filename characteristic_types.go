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
	DB.AddTableWithName(CharacteristicTypeBase{}, "characteristic_types").SetKeys(true, "Id")
}

type CharacteristicTypeService struct{}

// A CharacteristicType is a lookup type
type CharacteristicTypeBase struct {
	Id                     int64     `json:"id,omitempty"`
	CharacteristicTypeName string    `db:"characteristic_type_name" json:"characteristicTypeName"`
	CreatedAt              time.Time `db:"created_at" json:"createdAt"`
	UpdatedAt              time.Time `db:"updated_at" json:"updatedAt"`
	DeletedAt              NullTime  `db:"deleted_at" json:"deletedAt"`
	CreatedBy              int64     `db:"created_by" json:"createdBy"`
	UpdatedBy              int64     `db:"updated_by" json:"updatedBy"`
	DeletedBy              NullInt64 `db:"deleted_by" json:"deletedBy"`
}

type CharacteristicType struct {
	*CharacteristicTypeBase
	Characteristics NullSliceInt64 `db:"characteristics" json:"characteristics"`
	SortOrder       int64          `db:"sort_order" json:"sortOrder"`
}

type CharacteristicTypes []*CharacteristicType

type CharacteristicTypeJSON struct {
	CharacteristicType *CharacteristicType `json:"characteristicType"`
}

type CharacteristicTypesJSON struct {
	CharacteristicTypes *CharacteristicTypes `json:"characteristicTypes"`
}

func (c *CharacteristicType) marshal() ([]byte, error) {
	return json.Marshal(&CharacteristicTypeJSON{CharacteristicType: c})
}

func (c *CharacteristicTypes) marshal() ([]byte, error) {
	return json.Marshal(&CharacteristicTypesJSON{CharacteristicTypes: c})
}

func (c CharacteristicTypeService) list(val *url.Values) (entity, error) {
	if val == nil {
		return nil, errors.New("must provide options")
	}
	var opt ListOptions
	if err := schemaDecoder.Decode(&opt, *val); err != nil {
		return nil, err
	}

	var vals []interface{}
	sql := `SELECT ct.*, array_agg(c.id) AS characteristics,
			rank() OVER (ORDER BY ct.characteristic_type_name) AS sort_order
			FROM characteristic_types ct
			INNER JOIN characteristics c ON c.characteristic_type_id=ct.id`

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

	sql += " GROUP BY ct.id;"

	characteristic_types := make(CharacteristicTypes, 0)
	err := DBH.Select(&characteristic_types, sql, vals...)
	if err != nil {
		return nil, err
	}
	return &characteristic_types, nil
}

func (c CharacteristicTypeService) get(id int64, dummy string) (entity, error) {
	var characteristic_type CharacteristicType
	sql := `SELECT ct.*, array_agg(c.id) AS characteristics, 0 AS sort_order
			FROM characteristic_types ct
			INNER JOIN characteristics c ON c.characteristic_type_id=ct.id
			WHERE ct.id=$1
			GROUP BY ct.id;`
	if err := DBH.SelectOne(&characteristic_type, sql, id); err != nil {
		return nil, err
	}
	return &characteristic_type, nil
}
