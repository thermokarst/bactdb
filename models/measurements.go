package models

import (
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/thermokarst/bactdb/Godeps/_workspace/src/github.com/jmoiron/modl"
	"github.com/thermokarst/bactdb/errors"
	"github.com/thermokarst/bactdb/helpers"
	"github.com/thermokarst/bactdb/types"
)

func init() {
	DB.AddTableWithName(MeasurementBase{}, "measurements").SetKeys(true, "ID")
}

// PreInsert is a modl hook.
func (m *MeasurementBase) PreInsert(e modl.SqlExecutor) error {
	ct := helpers.CurrentTime()
	m.CreatedAt = ct
	m.UpdatedAt = ct
	return nil
}

// PreUpdate is a modl hook.
func (m *MeasurementBase) PreUpdate(e modl.SqlExecutor) error {
	m.UpdatedAt = helpers.CurrentTime()
	return nil
}

// UpdateError satisfies base interface.
func (m *MeasurementBase) UpdateError() error {
	return errors.ErrMeasurementNotUpdated
}

// DeleteError satisfies base interface.
func (m *MeasurementBase) DeleteError() error {
	return errors.ErrMeasurementNotDeleted
}

// MeasurementBase is what the DB expects for write operations
// There are three types of supported measurements: fixed-text, free-text,
// & numerical. The table has a constraint that will allow at most one
// for a particular combination of strain & characteristic.
type MeasurementBase struct {
	ID                    int64             `json:"id,omitempty"`
	StrainID              int64             `db:"strain_id" json:"strain"`
	CharacteristicID      int64             `db:"characteristic_id" json:"characteristic"`
	TextMeasurementTypeID types.NullInt64   `db:"text_measurement_type_id" json:"-"`
	TxtValue              types.NullString  `db:"txt_value" json:"-"`
	NumValue              types.NullFloat64 `db:"num_value" json:"-"`
	ConfidenceInterval    types.NullFloat64 `db:"confidence_interval" json:"confidenceInterval"`
	UnitTypeID            types.NullInt64   `db:"unit_type_id" json:"-"`
	Notes                 types.NullString  `db:"notes" json:"notes"`
	TestMethodID          types.NullInt64   `db:"test_method_id" json:"-"`
	CreatedAt             types.NullTime    `db:"created_at" json:"createdAt"`
	UpdatedAt             types.NullTime    `db:"updated_at" json:"updatedAt"`
	CreatedBy             int64             `db:"created_by" json:"createdBy"`
	UpdatedBy             int64             `db:"updated_by" json:"updatedBy"`
}

// Measurement is what the DB expects for read operations, and is what the API
// expects to return to the requester.
type Measurement struct {
	*MeasurementBase
	TextMeasurementType types.NullString `db:"text_measurement_type_name" json:"-"`
	UnitType            types.NullString `db:"unit_type_name" json:"unitType"`
	TestMethod          types.NullString `db:"test_method_name" json:"testMethod"`
	CanEdit             bool             `db:"-" json:"canEdit"`
}

// FakeMeasurement is a dummy struct to prevent infinite-loop/stack overflow on serialization.
type FakeMeasurement Measurement

// MarshalJSON is custom JSON serialization to handle multi-type "Value".
func (m *Measurement) MarshalJSON() ([]byte, error) {
	fm := FakeMeasurement(*m)
	return json.Marshal(struct {
		*FakeMeasurement
		Value string `json:"value"`
	}{
		FakeMeasurement: &fm,
		Value:           m.Value(),
	})
}

// UnmarshalJSON is custom JSON deserialization to handle multi-type "Value"
func (m *Measurement) UnmarshalJSON(b []byte) error {
	var measurement struct {
		FakeMeasurement
		Value interface{} `json:"value"`
	}
	if err := json.Unmarshal(b, &measurement); err != nil {
		return err
	}

	switch v := measurement.Value.(type) {
	case string:
		// Test if actually a lookup
		id, err := GetTextMeasurementTypeID(v)
		if err != nil {
			if err == sql.ErrNoRows {
				measurement.TxtValue = types.NullString{sql.NullString{String: v, Valid: true}}
			} else {
				return err
			}
		} else {
			measurement.TextMeasurementTypeID = types.NullInt64{sql.NullInt64{Int64: id, Valid: true}}
		}
	case int64:
		measurement.NumValue = types.NullFloat64{sql.NullFloat64{Float64: float64(v), Valid: true}}
	case float64:
		measurement.NumValue = types.NullFloat64{sql.NullFloat64{Float64: v, Valid: true}}
	}

	*m = Measurement(measurement.FakeMeasurement)

	return nil
}

// Value returns the value of the measurement
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

// Measurements are multiple measurement entities
type Measurements []*Measurement

// MeasurementMeta stashes some metadata related to the entity
type MeasurementMeta struct {
	CanAdd bool `json:"canAdd"`
}

// ListMeasurements returns all measurements
func ListMeasurements(opt helpers.MeasurementListOptions, claims *types.Claims) (*Measurements, error) {
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

	strainIDs := len(opt.Strains) != 0
	charIDs := len(opt.Characteristics) != 0
	ids := len(opt.IDs) != 0

	if strainIDs || charIDs || ids {
		var paramsCounter int64 = 2
		q += "\nWHERE ("

		// Filter by strains
		if strainIDs {
			q += helpers.ValsIn("st.id", opt.Strains, &vals, &paramsCounter)
		}

		if strainIDs && (charIDs || ids) {
			q += " AND "
		}

		// Filter by characteristics
		if charIDs {
			q += helpers.ValsIn("c.id", opt.Characteristics, &vals, &paramsCounter)
		}

		if charIDs && ids {
			q += " AND "
		}

		// Get specific records
		if ids {
			q += helpers.ValsIn("m.id", opt.IDs, &vals, &paramsCounter)
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
		m.CanEdit = helpers.CanEdit(claims, m.CreatedBy)
	}

	return &measurements, nil
}

// GetMeasurement returns a particular measurement.
func GetMeasurement(id int64, genus string, claims *types.Claims) (*Measurement, error) {
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
			return nil, errors.ErrMeasurementNotFound
		}
		return nil, err
	}

	measurement.CanEdit = helpers.CanEdit(claims, measurement.CreatedBy)

	return &measurement, nil
}

// CharacteristicOptsFromMeasurements returns the options for finding all related
// characteristics for a set of measurements.
func CharacteristicOptsFromMeasurements(opt helpers.MeasurementListOptions) (*helpers.ListOptions, error) {
	return &helpers.ListOptions{Genus: opt.Genus, IDs: opt.Characteristics}, nil
}

// StrainOptsFromMeasurements returns the options for finding all related
// strains from a set of measurements.
func StrainOptsFromMeasurements(opt helpers.MeasurementListOptions) (*helpers.ListOptions, error) {
	return &helpers.ListOptions{Genus: opt.Genus, IDs: opt.Strains}, nil
}

// GetTextMeasurementTypeID returns the ID for a particular text measurement type
func GetTextMeasurementTypeID(val string) (int64, error) {
	var id int64
	q := `SELECT id FROM text_measurement_types WHERE text_measurement_name=$1;`

	if err := DBH.SelectOne(&id, q, val); err != nil {
		return 0, err
	}
	return id, nil
}
