package models

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/thermokarst/bactdb/Godeps/_workspace/src/github.com/jmoiron/modl"
	"github.com/thermokarst/bactdb/helpers"
	"github.com/thermokarst/bactdb/types"
)

var (
	ErrCharacteristicNotFound   = errors.New("Characteristic not found")
	ErrCharacteristicNotUpdated = errors.New("Characteristic not updated")
)

func init() {
	DB.AddTableWithName(CharacteristicBase{}, "characteristics").SetKeys(true, "Id")
}

func (c *CharacteristicBase) PreInsert(e modl.SqlExecutor) error {
	ct := helpers.CurrentTime()
	c.CreatedAt = ct
	c.UpdatedAt = ct
	return nil
}

func (c *CharacteristicBase) PreUpdate(e modl.SqlExecutor) error {
	c.UpdatedAt = helpers.CurrentTime()
	return nil
}

type CharacteristicBase struct {
	Id                   int64           `json:"id,omitempty"`
	CharacteristicName   string          `db:"characteristic_name" json:"characteristicName"`
	CharacteristicTypeId int64           `db:"characteristic_type_id" json:"-"`
	SortOrder            types.NullInt64 `db:"sort_order" json:"sortOrder"`
	CreatedAt            types.NullTime  `db:"created_at" json:"createdAt"`
	UpdatedAt            types.NullTime  `db:"updated_at" json:"updatedAt"`
	DeletedAt            types.NullTime  `db:"deleted_at" json:"deletedAt"`
	CreatedBy            int64           `db:"created_by" json:"createdBy"`
	UpdatedBy            int64           `db:"updated_by" json:"updatedBy"`
	DeletedBy            types.NullInt64 `db:"deleted_by" json:"deletedBy"`
}

type Characteristic struct {
	*CharacteristicBase
	Measurements       types.NullSliceInt64 `db:"measurements" json:"measurements"`
	Strains            types.NullSliceInt64 `db:"strains" json:"strains"`
	CharacteristicType string               `db:"characteristic_type_name" json:"characteristicTypeName"`
	CanEdit            bool                 `db:"-" json:"canEdit"`
}

type Characteristics []*Characteristic

type CharacteristicMeta struct {
	CanAdd bool `json:"canAdd"`
}

func ListCharacteristics(opt helpers.ListOptions, claims *types.Claims) (*Characteristics, error) {
	var vals []interface{}

	q := `SELECT c.*, ct.characteristic_type_name,
			array_agg(DISTINCT st.id) AS strains, array_agg(DISTINCT m.id) AS measurements
			FROM strains st
			INNER JOIN species sp ON sp.id=st.species_id
			INNER JOIN genera g ON g.id=sp.genus_id AND LOWER(g.genus_name)=LOWER($1)
			INNER JOIN measurements m ON m.strain_id=st.id
			RIGHT OUTER JOIN characteristics c ON c.id=m.characteristic_id
			INNER JOIN characteristic_types ct ON ct.id=c.characteristic_type_id`
	vals = append(vals, opt.Genus)

	if len(opt.Ids) != 0 {
		var counter int64 = 2
		w := helpers.ValsIn("c.id", opt.Ids, &vals, &counter)

		q += fmt.Sprintf(" WHERE %s", w)
	}

	q += ` GROUP BY c.id, ct.characteristic_type_name
			ORDER BY ct.characteristic_type_name, c.sort_order ASC;`

	var characteristics Characteristics
	err := DBH.Select(&characteristics, q, vals...)
	if err != nil {
		return nil, err
	}

	for _, c := range characteristics {
		c.CanEdit = helpers.CanEdit(claims, c.CreatedBy)
	}

	return &characteristics, nil
}

func StrainOptsFromCharacteristics(opt helpers.ListOptions) (*helpers.ListOptions, error) {
	relatedStrainIds := make([]int64, 0)
	baseQ := `SELECT DISTINCT m.strain_id
		FROM measurements m
		INNER JOIN strains st ON st.id=m.strain_id
		INNER JOIN species sp ON sp.id=st.species_id
		INNER JOIN genera g ON g.id=sp.genus_id AND LOWER(g.genus_name)=LOWER($1)`
	if opt.Ids == nil {
		q := fmt.Sprintf("%s;", baseQ)
		if err := DBH.Select(&relatedStrainIds, q, opt.Genus); err != nil {
			return nil, err
		}
	} else {
		var vals []interface{}
		var count int64 = 2
		vals = append(vals, opt.Genus)
		q := fmt.Sprintf("%s WHERE %s ", baseQ, helpers.ValsIn("m.characteristic_id", opt.Ids, &vals, &count))

		if err := DBH.Select(&relatedStrainIds, q, vals...); err != nil {
			return nil, err
		}
	}

	return &helpers.ListOptions{Genus: opt.Genus, Ids: relatedStrainIds}, nil
}

func MeasurementOptsFromCharacteristics(opt helpers.ListOptions) (*helpers.MeasurementListOptions, error) {
	relatedMeasurementIds := make([]int64, 0)
	baseQ := `SELECT m.id
		FROM measurements m
		INNER JOIN strains st ON st.id=m.strain_id
		INNER JOIN species sp ON sp.id=st.species_id
		INNER JOIN genera g ON g.id=sp.genus_id AND LOWER(g.genus_name)=LOWER($1)`

	if opt.Ids == nil {
		q := fmt.Sprintf("%s;", baseQ)
		if err := DBH.Select(&relatedMeasurementIds, q, opt.Genus); err != nil {
			return nil, err
		}
	} else {
		var vals []interface{}
		var count int64 = 2
		vals = append(vals, opt.Genus)
		q := fmt.Sprintf("%s WHERE %s;", baseQ, helpers.ValsIn("characteristic_id", opt.Ids, &vals, &count))

		if err := DBH.Select(&relatedMeasurementIds, q, vals...); err != nil {
			return nil, err
		}
	}

	return &helpers.MeasurementListOptions{ListOptions: helpers.ListOptions{Genus: opt.Genus, Ids: relatedMeasurementIds}, Strains: nil, Characteristics: nil}, nil
}

func StrainsFromCharacteristicId(id int64, genus string, claims *types.Claims) (*Strains, *helpers.ListOptions, error) {
	opt := helpers.ListOptions{
		Genus: genus,
		Ids:   []int64{id},
	}

	strains_opt, err := StrainOptsFromCharacteristics(opt)
	if err != nil {
		return nil, nil, err
	}

	strains, err := ListStrains(*strains_opt, claims)
	if err != nil {
		return nil, nil, err
	}

	return strains, strains_opt, nil
}

func MeasurementsFromCharacteristicId(id int64, genus string, claims *types.Claims) (*Measurements, *helpers.MeasurementListOptions, error) {
	opt := helpers.ListOptions{
		Genus: genus,
		Ids:   []int64{id},
	}

	measurement_opt, err := MeasurementOptsFromCharacteristics(opt)
	if err != nil {
		return nil, nil, err
	}

	measurements, err := ListMeasurements(*measurement_opt, claims)
	if err != nil {
		return nil, nil, err
	}

	return measurements, measurement_opt, nil
}

func GetCharacteristic(id int64, genus string, claims *types.Claims) (*Characteristic, error) {
	var characteristic Characteristic
	q := `SELECT c.*, ct.characteristic_type_name,
		array_agg(DISTINCT st.id) AS strains, array_agg(DISTINCT m.id) AS measurements
		FROM strains st
		INNER JOIN species sp ON sp.id=st.species_id
		INNER JOIN genera g ON g.id=sp.genus_id AND LOWER(g.genus_name)=LOWER($1)
		INNER JOIN measurements m ON m.strain_id=st.id
		RIGHT OUTER JOIN characteristics c ON c.id=m.characteristic_id
		INNER JOIN characteristic_types ct ON ct.id=c.characteristic_type_id
		WHERE c.id=$2
		GROUP BY c.id, ct.characteristic_type_name;`
	if err := DBH.SelectOne(&characteristic, q, genus, id); err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrCharacteristicNotFound
		}
		return nil, err
	}

	characteristic.CanEdit = helpers.CanEdit(claims, characteristic.CreatedBy)

	return &characteristic, nil
}

func InsertOrGetCharacteristicType(val string, claims *types.Claims) (int64, error) {
	var id int64
	q := `SELECT id FROM characteristic_types WHERE characteristic_type_name=$1;`
	if err := DBH.SelectOne(&id, q, val); err != nil {
		if err == sql.ErrNoRows {
			i := `INSERT INTO characteristic_types
				(characteristic_type_name, created_at, updated_at, created_by, updated_by)
				VALUES ($1, $2, $3, $4, $5) RETURNING id;`
			ct := helpers.CurrentTime()
			var result sql.Result
			var insertErr error
			stmt, err := DB.Db.Prepare(i)
			if result, insertErr = stmt.Exec(val, ct, ct, claims.Sub, claims.Sub); insertErr != nil {
				return 0, insertErr
			}
			id, err = result.LastInsertId()
			if err != nil {
				return 0, err
			}
		} else {
			return 0, err
		}
	}
	return id, nil
}
