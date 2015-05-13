package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"
)

func init() {
	DB.AddTableWithName(CharacteristicBase{}, "characteristics").SetKeys(true, "Id")
}

// A Characteristic is a lookup type
type CharacteristicBase struct {
	Id                   int64     `json:"id,omitempty"`
	CharacteristicName   string    `db:"characteristic_name" json:"characteristicName"`
	CharacteristicTypeId int64     `db:"characteristic_type_id" json:"-"`
	CreatedAt            time.Time `db:"created_at" json:"createdAt"`
	UpdatedAt            time.Time `db:"updated_at" json:"updatedAt"`
	DeletedAt            NullTime  `db:"deleted_at" json:"deletedAt"`
	CreatedBy            int64     `db:"created_by" json:"createdBy"`
	UpdatedBy            int64     `db:"updated_by" json:"updatedBy"`
	DeletedBy            NullInt64 `db:"deleted_by" json:"deletedBy"`
}

type Characteristic struct {
	*CharacteristicBase
	Measurements           NullSliceInt64 `db:"measurements" json:"measurements"`
	Strains                NullSliceInt64 `db:"strains" json:"strains"`
	CharacteristicTypeName string         `db:"characteristic_type_name" json:"characteristicType"`
}

type CharacteristicJSON struct {
	Characteristic *Characteristic `json:"characteristic"`
}

type CharacteristicsJSON struct {
	Characteristics []*Characteristic `json:"characteristics"`
}

type CharacteristicListOptions struct {
	ListOptions
	Genus string
}

func serveCharacteristicsList(w http.ResponseWriter, r *http.Request) {
	var opt CharacteristicListOptions
	if err := schemaDecoder.Decode(&opt, r.URL.Query()); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	opt.Genus = mux.Vars(r)["genus"]

	characteristics, err := dbGetCharacteristics(&opt)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if characteristics == nil {
		characteristics = []*Characteristic{}
	}
	data, err := json.Marshal(CharacteristicsJSON{Characteristics: characteristics})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Write(data)
}

func serveCharacteristic(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(mux.Vars(r)["Id"], 10, 0)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	characteristic, err := dbGetCharacteristic(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	data, err := json.Marshal(CharacteristicJSON{Characteristic: characteristic})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Write(data)
}

func dbGetCharacteristics(opt *CharacteristicListOptions) ([]*Characteristic, error) {
	if opt == nil {
		return nil, errors.New("must provide options")
	}
	var vals []interface{}
	sql := `SELECT c.*, ct.characteristic_type_name,
			array_agg(m.id) AS measurements, array_agg(st.id) AS strains
			FROM characteristics c
			INNER JOIN characteristic_types ct ON ct.id=characteristic_type_id
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

	sql += " GROUP BY c.id, ct.characteristic_type_name;"

	var characteristics []*Characteristic
	err := DBH.Select(&characteristics, sql, vals...)
	if err != nil {
		return nil, err
	}
	return characteristics, nil
}

func dbGetCharacteristic(id int64) (*Characteristic, error) {
	var characteristic Characteristic
	sql := `SELECT c.*, ct.characteristic_type_name,
			array_agg(m.id) AS measurements, array_agg(st.id) AS strains
			FROM characteristics c
			INNER JOIN characteristic_types ct ON ct.id=characteristic_type_id
			LEFT OUTER JOIN measurements m ON m.characteristic_id=c.id
			LEFT OUTER JOIN strains st ON st.id=m.strain_id
			WHERE c.id=$1
			GROUP BY c.id, ct.characteristic_type_name;`
	if err := DBH.SelectOne(&characteristic, sql, id); err != nil {
		return nil, err
	}
	return &characteristic, nil
}
