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

var ErrMeasurementNotFound = errors.New("measurement not found")

func init() {
	DB.AddTableWithName(MeasurementBase{}, "measurements").SetKeys(true, "Id")
}

// There are three types of supported measurements: fixed-test, free-text,
// & numerical. The table has a constraint that will allow one or the other
// for a particular combination of strain & characteristic, but not both.
// MeasurementBase is what the DB expects to see for inserts/updates
type MeasurementBase struct {
	Id                    int64       `json:"id,omitempty"`
	StrainId              int64       `db:"strain_id" json:"strain"`
	CharacteristicId      int64       `db:"characteristic_id" json:"-"`
	TextMeasurementTypeId NullInt64   `db:"text_measurement_type_id" json:"-"`
	TxtValue              NullString  `db:"txt_value" json:"txtValue"`
	NumValue              NullFloat64 `db:"num_value" json:"numValue"`
	ConfidenceInterval    NullFloat64 `db:"confidence_interval" json:"confidenceInterval"`
	UnitTypeId            NullInt64   `db:"unit_type_id" json:"-"`
	Notes                 NullString  `db:"notes" json:"notes"`
	TestMethodId          NullInt64   `db:"test_method_id" json:"-"`
	CreatedAt             time.Time   `db:"created_at" json:"createdAt"`
	UpdatedAt             time.Time   `db:"updated_at" json:"updatedAt"`
	CreatedBy             int64       `db:"created_by" json:"createdBy"`
	UpdatedBy             int64       `db:"updated_by" json:"updatedBy"`
}

// Measurement & MeasurementJSON(s) are what ember expects to see
type Measurement struct {
	*MeasurementBase
	Characteristic      NullString `db:"characteristic_name" json:"characteristic"`
	TextMeasurementType NullString `db:"text_measurement_type_name" json:"textMeasurementType"`
	UnitType            NullString `db:"unit_type_name" json:"unitType"`
	TestMethod          NullString `db:"test_method_name" json:"testMethod"`
}

type MeasurementJSON struct {
	Measurement *Measurement `json:"measurement"`
}

type MeasurementsJSON struct {
	Measurements []*Measurement `json:"measurements"`
}

type MeasurementListOptions struct {
	ListOptions
	Genus string
}

func serveMeasurementsList(w http.ResponseWriter, r *http.Request) {
	var opt MeasurementListOptions
	if err := schemaDecoder.Decode(&opt, r.URL.Query()); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	opt.Genus = mux.Vars(r)["genus"]

	measurements, err := dbGetMeasurements(&opt)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if measurements == nil {
		measurements = []*Measurement{}
	}
	data, err := json.Marshal(MeasurementsJSON{Measurements: measurements})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Write(data)
}

func serveMeasurement(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(mux.Vars(r)["Id"], 10, 0)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	measurement, err := dbGetMeasurement(id, mux.Vars(r)["genus"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	data, err := json.Marshal(MeasurementJSON{Measurement: measurement})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Write(data)
}

func dbGetMeasurements(opt *MeasurementListOptions) ([]*Measurement, error) {
	if opt == nil {
		return nil, errors.New("must provide options")
	}

	var vals []interface{}
	sql := `SELECT m.*, c.characteristic_name,
		t.text_measurement_name AS text_measurement_type_name,
		u.symbol AS unit_type_name, te.name AS test_method_name
		FROM measurements m
		INNER JOIN strains st ON st.id=m.strain_id
		INNER JOIN species sp ON sp.id=st.species_id
		INNER JOIN genera g ON g.id=sp.genus_id AND LOWER(g.genus_name)=$1
		LEFT OUTER JOIN characteristics c ON c.id=m.characteristic_id
		LEFT OUTER JOIN text_measurement_types t ON t.id=m.text_measurement_type_id
		LEFT OUTER JOIN unit_types u ON u.id=m.unit_type_id
		LEFT OUTER JOIN test_methods te ON te.id=m.test_method_id`
	vals = append(vals, opt.Genus)

	if len(opt.Ids) != 0 {
		var conds []string

		m := "m.id IN ("
		for i, id := range opt.Ids {
			m = m + fmt.Sprintf("$%v,", i+2) // start param index at 2
			vals = append(vals, id)
		}
		m = m[:len(m)-1] + ")"
		conds = append(conds, m)
		sql += " WHERE (" + strings.Join(conds, ") AND (") + ")"
	}

	sql += ";"

	var measurements []*Measurement
	err := DBH.Select(&measurements, sql, vals...)
	if err != nil {
		return nil, err
	}
	return measurements, nil
}

func dbGetMeasurement(id int64, genus string) (*Measurement, error) {
	var measurement Measurement
	sql := `SELECT m.*, c.characteristic_name,
		t.text_measurement_name AS text_measurement_type_name,
		u.symbol AS unit_type_name, te.name AS test_method_name
		FROM measurements m
		INNER JOIN strains st ON st.id=m.strain_id
		INNER JOIN species sp ON sp.id=st.species_id
		INNER JOIN genera g ON g.id=sp.genus_id AND LOWER(g.genus_name)=$1
		LEFT OUTER JOIN characteristics c ON c.id=m.characteristic_id
		LEFT OUTER JOIN text_measurement_types t ON t.id=m.text_measurement_type_id
		LEFT OUTER JOIN unit_types u ON u.id=m.unit_type_id
		LEFT OUTER JOIN test_methods te ON te.id=m.test_method_id
		WHERE m.id=$2;`
	if err := DBH.SelectOne(&measurement, sql, genus, id); err != nil {
		return nil, err
	}
	if &measurement == nil {
		return nil, ErrMeasurementNotFound
	}
	return &measurement, nil
}
