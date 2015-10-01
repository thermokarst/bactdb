package models

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/thermokarst/bactdb/Godeps/_workspace/src/github.com/jmoiron/modl"
	"github.com/thermokarst/bactdb/helpers"
	"github.com/thermokarst/bactdb/types"
)

var (
	ErrMeasurementNotFound     = errors.New("Measurement not found")
	ErrMeasurementNotFoundJSON = types.NewJSONError(ErrMeasurementNotFound, http.StatusNotFound)
)

func init() {
	DB.AddTableWithName(MeasurementBase{}, "measurements").SetKeys(true, "Id")
}

func (m *MeasurementBase) PreInsert(e modl.SqlExecutor) error {
	ct := helpers.CurrentTime()
	m.CreatedAt = ct
	m.UpdatedAt = ct
	return nil
}

func (m *MeasurementBase) PreUpdate(e modl.SqlExecutor) error {
	m.UpdatedAt = helpers.CurrentTime()
	return nil
}

// There are three types of supported measurements: fixed-text, free-text,
// & numerical. The table has a constraint that will allow at most one
// for a particular combination of strain & characteristic.
// MeasurementBase is what the DB expects to see for inserts/updates
type MeasurementBase struct {
	Id                    int64             `json:"id,omitempty"`
	StrainId              int64             `db:"strain_id" json:"strain"`
	CharacteristicId      int64             `db:"characteristic_id" json:"characteristic"`
	TextMeasurementTypeId types.NullInt64   `db:"text_measurement_type_id" json:"-"`
	TxtValue              types.NullString  `db:"txt_value" json:"-"`
	NumValue              types.NullFloat64 `db:"num_value" json:"-"`
	ConfidenceInterval    types.NullFloat64 `db:"confidence_interval" json:"confidenceInterval"`
	UnitTypeId            types.NullInt64   `db:"unit_type_id" json:"-"`
	Notes                 types.NullString  `db:"notes" json:"notes"`
	TestMethodId          types.NullInt64   `db:"test_method_id" json:"-"`
	CreatedAt             types.NullTime    `db:"created_at" json:"createdAt"`
	UpdatedAt             types.NullTime    `db:"updated_at" json:"updatedAt"`
	CreatedBy             int64             `db:"created_by" json:"createdBy"`
	UpdatedBy             int64             `db:"updated_by" json:"updatedBy"`
}

type Measurement struct {
	*MeasurementBase
	TextMeasurementType types.NullString `db:"text_measurement_type_name" json:"-"`
	UnitType            types.NullString `db:"unit_type_name" json:"unitType"`
	TestMethod          types.NullString `db:"test_method_name" json:"testMethod"`
	CanEdit             bool             `db:"-" json:"canEdit"`
}

type FakeMeasurement Measurement

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
		id, err := GetTextMeasurementTypeId(v)
		if err != nil {
			if err == sql.ErrNoRows {
				measurement.TxtValue = types.NullString{sql.NullString{String: v, Valid: true}}
			} else {
				return err
			}
		} else {
			measurement.TextMeasurementTypeId = types.NullInt64{sql.NullInt64{Int64: id, Valid: true}}
		}
	case int64:
		measurement.NumValue = types.NullFloat64{sql.NullFloat64{Float64: float64(v), Valid: true}}
	case float64:
		measurement.NumValue = types.NullFloat64{sql.NullFloat64{Float64: v, Valid: true}}
	}

	*m = Measurement(measurement.FakeMeasurement)

	return nil
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

	strainIds := len(opt.Strains) != 0
	charIds := len(opt.Characteristics) != 0
	ids := len(opt.Ids) != 0

	if strainIds || charIds || ids {
		var paramsCounter int64 = 2
		q += "\nWHERE ("

		// Filter by strains
		if strainIds {
			q += helpers.ValsIn("st.id", opt.Strains, &vals, &paramsCounter)
		}

		if strainIds && (charIds || ids) {
			q += " AND "
		}

		// Filter by characteristics
		if charIds {
			q += helpers.ValsIn("c.id", opt.Characteristics, &vals, &paramsCounter)
		}

		if charIds && ids {
			q += " AND "
		}

		// Get specific records
		if ids {
			q += helpers.ValsIn("m.id", opt.Ids, &vals, &paramsCounter)
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
			return nil, ErrMeasurementNotFound
		}
		return nil, err
	}

	measurement.CanEdit = helpers.CanEdit(claims, measurement.CreatedBy)

	return &measurement, nil
}

func CharacteristicOptsFromMeasurements(opt helpers.MeasurementListOptions) (*helpers.ListOptions, error) {
	return &helpers.ListOptions{Genus: opt.Genus, Ids: opt.Characteristics}, nil
}

func StrainOptsFromMeasurements(opt helpers.MeasurementListOptions) (*helpers.ListOptions, error) {
	return &helpers.ListOptions{Genus: opt.Genus, Ids: opt.Strains}, nil
}

func GetTextMeasurementTypeId(val string) (int64, error) {
	var id int64
	q := `SELECT id FROM text_measurement_types WHERE text_measurement_name=$1;`

	if err := DBH.SelectOne(&id, q, val); err != nil {
		return 0, err
	}
	return id, nil
}
