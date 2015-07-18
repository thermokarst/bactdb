package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"

	"github.com/jmoiron/modl"
)

var (
	ErrMeasurementNotFound     = errors.New("Measurement not found")
	ErrMeasurementNotFoundJSON = newJSONError(ErrMeasurementNotFound, http.StatusNotFound)
)

func init() {
	DB.AddTableWithName(MeasurementBase{}, "measurements").SetKeys(true, "Id")
}

func (m *MeasurementBase) PreInsert(e modl.SqlExecutor) error {
	ct := currentTime()
	m.CreatedAt = ct
	m.UpdatedAt = ct
	return nil
}

func (m *MeasurementBase) PreUpdate(e modl.SqlExecutor) error {
	m.UpdatedAt = currentTime()
	return nil
}

type MeasurementService struct{}

// There are three types of supported measurements: fixed-test, free-text,
// & numerical. The table has a constraint that will allow one or the other
// for a particular combination of strain & characteristic, but not both.
// MeasurementBase is what the DB expects to see for inserts/updates
type MeasurementBase struct {
	Id                    int64       `json:"id,omitempty"`
	StrainId              int64       `db:"strain_id" json:"strain"`
	CharacteristicId      int64       `db:"characteristic_id" json:"characteristic"`
	TextMeasurementTypeId NullInt64   `db:"text_measurement_type_id" json:"-"`
	TxtValue              NullString  `db:"txt_value" json:"txtValue"`
	NumValue              NullFloat64 `db:"num_value" json:"numValue"`
	ConfidenceInterval    NullFloat64 `db:"confidence_interval" json:"confidenceInterval"`
	UnitTypeId            NullInt64   `db:"unit_type_id" json:"-"`
	Notes                 NullString  `db:"notes" json:"notes"`
	TestMethodId          NullInt64   `db:"test_method_id" json:"-"`
	CreatedAt             NullTime    `db:"created_at" json:"createdAt"`
	UpdatedAt             NullTime    `db:"updated_at" json:"updatedAt"`
	CreatedBy             int64       `db:"created_by" json:"createdBy"`
	UpdatedBy             int64       `db:"updated_by" json:"updatedBy"`
}

type Measurement struct {
	*MeasurementBase
	TextMeasurementType NullString `db:"text_measurement_type_name" json:"textMeasurementType"`
	UnitType            NullString `db:"unit_type_name" json:"unitType"`
	TestMethod          NullString `db:"test_method_name" json:"testMethod"`
	CanEdit             bool       `db:"-" json:"canEdit"`
}

func (m *Measurement) Value() string {
	if m.TextMeasurementType.Valid {
		return m.TextMeasurementType.String
	}
	if m.TxtValue.Valid {
		return m.TxtValue.String
	}
	if m.NumValue.Valid {
		return fmt.Sprintf("%f", m.NumValue.Float64)
	}
	return ""
}

type Measurements []*Measurement

type MeasurementMeta struct {
	CanAdd bool `json:"canAdd"`
}

type MeasurementPayload struct {
	Measurement *Measurement `json:"measurement"`
}

type MeasurementsPayload struct {
	Strains         *Strains         `json:"strains"`
	Characteristics *Characteristics `json:"characteristics"`
	Measurements    *Measurements    `json:"measurements"`
}

func (m *MeasurementPayload) marshal() ([]byte, error) {
	return json.Marshal(m)
}

func (m *MeasurementsPayload) marshal() ([]byte, error) {
	return json.Marshal(m)
}

type MeasurementListOptions struct {
	ListOptions
	Strains         []int64 `schema:"strain_ids"`
	Characteristics []int64 `schema:"characteristic_ids"`
}

func (m MeasurementService) list(val *url.Values, claims *Claims) (entity, *appError) {
	if val == nil {
		return nil, ErrMustProvideOptionsJSON
	}
	var opt MeasurementListOptions
	if err := schemaDecoder.Decode(&opt, *val); err != nil {
		return nil, newJSONError(err, http.StatusInternalServerError)
	}

	measurements, err := listMeasurements(opt, claims)
	if err != nil {
		return nil, newJSONError(err, http.StatusInternalServerError)
	}

	char_opts, err := characteristicOptsFromMeasurements(opt)
	if err != nil {
		return nil, newJSONError(err, http.StatusInternalServerError)
	}

	characteristics, err := listCharacteristics(*char_opts, claims)
	if err != nil {
		return nil, newJSONError(err, http.StatusInternalServerError)
	}

	strain_opts, err := strainOptsFromMeasurements(opt)
	if err != nil {
		return nil, newJSONError(err, http.StatusInternalServerError)
	}

	strains, err := listStrains(*strain_opts, claims)
	if err != nil {
		return nil, newJSONError(err, http.StatusInternalServerError)
	}

	payload := MeasurementsPayload{
		Characteristics: characteristics,
		Strains:         strains,
		Measurements:    measurements,
	}

	return &payload, nil
}

func (m MeasurementService) get(id int64, genus string, claims *Claims) (entity, *appError) {
	measurement, err := getMeasurement(id, genus, claims)
	if err != nil {
		return nil, newJSONError(err, http.StatusInternalServerError)
	}

	payload := MeasurementPayload{
		Measurement: measurement,
	}

	return &payload, nil
}

func listMeasurements(opt MeasurementListOptions, claims *Claims) (*Measurements, error) {
	var vals []interface{}

	q := `SELECT m.*, t.text_measurement_name AS text_measurement_type_name,
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

	strainIds := len(opt.Strains) != 0
	charIds := len(opt.Characteristics) != 0
	ids := len(opt.Ids) != 0

	if strainIds || charIds || ids {
		var paramsCounter int64 = 2
		q += "\nWHERE ("

		// Filter by strains
		if strainIds {
			q += valsIn("st.id", opt.Strains, &vals, &paramsCounter)
		}

		if strainIds && (charIds || ids) {
			q += " AND "
		}

		// Filter by characteristics
		if charIds {
			q += valsIn("c.id", opt.Characteristics, &vals, &paramsCounter)
		}

		if charIds && ids {
			q += " AND "
		}

		// Get specific records
		if ids {
			q += valsIn("m.id", opt.Ids, &vals, &paramsCounter)
		}
		q += ")"
	}
	q += ";"

	measurements := make(Measurements, 0)
	err := DBH.Select(&measurements, q, vals...)
	if err != nil {
		return nil, err
	}

	for _, m := range measurements {
		m.CanEdit = canEdit(claims, m.CreatedBy)
	}

	return &measurements, nil
}

func getMeasurement(id int64, genus string, claims *Claims) (*Measurement, error) {
	var measurement Measurement

	q := `SELECT m.*, t.text_measurement_name AS text_measurement_type_name,
		u.symbol AS unit_type_name, te.name AS test_method_name
		FROM measurements m
		INNER JOIN strains st ON st.id=m.strain_id
		INNER JOIN species sp ON sp.id=st.species_id
		INNER JOIN genera g ON g.id=sp.genus_id AND LOWER(g.genus_name)=LOWER($1)
		LEFT OUTER JOIN characteristics c ON c.id=m.characteristic_id
		LEFT OUTER JOIN text_measurement_types t ON t.id=m.text_measurement_type_id
		LEFT OUTER JOIN unit_types u ON u.id=m.unit_type_id
		LEFT OUTER JOIN test_methods te ON te.id=m.test_method_id
		WHERE m.id=$2;`
	if err := DBH.SelectOne(&measurement, q, genus, id); err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrMeasurementNotFound
		}
		return nil, err
	}

	measurement.CanEdit = canEdit(claims, measurement.CreatedBy)

	return &measurement, nil
}

func characteristicOptsFromMeasurements(opt MeasurementListOptions) (*ListOptions, error) {
	return &ListOptions{Genus: opt.Genus, Ids: opt.Characteristics}, nil
}

func strainOptsFromMeasurements(opt MeasurementListOptions) (*ListOptions, error) {
	return &ListOptions{Genus: opt.Genus, Ids: opt.Strains}, nil
}
