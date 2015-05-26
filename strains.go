package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/context"
	"github.com/gorilla/mux"
)

var (
	ErrStrainNotFound   = errors.New("strain not found")
	ErrStrainNotUpdated = errors.New("strain not updated")
)

func init() {
	DB.AddTableWithName(StrainBase{}, "strains").SetKeys(true, "Id")
}

// StrainBase is what the DB expects to see for inserts/updates
type StrainBase struct {
	Id               int64      `db:"id" json:"id"`
	SpeciesId        int64      `db:"species_id" json:"-"`
	StrainName       string     `db:"strain_name" json:"strainName"`
	TypeStrain       bool       `db:"type_strain" json:"typeStrain"`
	AccessionNumbers string     `db:"accession_numbers" json:"accessionNumbers"`
	Genbank          NullString `db:"genbank" json:"genbank"`
	IsolatedFrom     NullString `db:"isolated_from" json:"isolatedFrom"`
	Notes            NullString `db:"notes" json:"notes"`
	CreatedAt        time.Time  `db:"created_at" json:"createdAt"`
	UpdatedAt        time.Time  `db:"updated_at" json:"updatedAt"`
	DeletedAt        NullTime   `db:"deleted_at" json:"deletedAt"`
	CreatedBy        int64      `db:"created_by" json:"createdBy"`
	UpdatedBy        int64      `db:"updated_by" json:"updatedBy"`
	DeletedBy        NullInt64  `db:"deleted_by" json:"deletedBy"`
}

// Strain & StrainJSON(s) are what ember expects to see
type Strain struct {
	*StrainBase
	SpeciesName       string         `db:"species_name" json:"speciesName"`
	Measurements      NullSliceInt64 `db:"measurements" json:"measurements"`
	TotalMeasurements int            `db:"total_measurements" json:"totalMeasurements"`
}

type StrainJSON struct {
	Strain *Strain `json:"strain"`
}

type StrainsJSON struct {
	Strains []*Strain `json:"strains"`
}

type StrainListOptions struct {
	ListOptions
	Genus string
}

func serveStrainsList(w http.ResponseWriter, r *http.Request) {
	var opt StrainListOptions
	if err := schemaDecoder.Decode(&opt, r.URL.Query()); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	opt.Genus = mux.Vars(r)["genus"]

	strains, err := dbGetStrains(&opt)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if strains == nil {
		strains = []*Strain{}
	}
	data, err := json.Marshal(StrainsJSON{Strains: strains})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Write(data)
}

func serveStrain(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(mux.Vars(r)["Id"], 10, 0)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	strain, err := dbGetStrain(id, mux.Vars(r)["genus"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data, err := json.Marshal(StrainJSON{Strain: strain})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Write(data)
}

func serveUpdateStrain(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(mux.Vars(r)["Id"], 10, 0)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var strainjson StrainJSON
	err = json.NewDecoder(r.Body).Decode(&strainjson)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	c := context.Get(r, "claims")
	var claims Claims = c.(Claims)
	strainjson.Strain.UpdatedBy = claims.Sub
	strainjson.Strain.UpdatedAt = time.Now()
	strainjson.Strain.Id = id

	err = dbUpdateStrain(strainjson.Strain)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var strain *Strain
	strain, err = dbGetStrain(id, mux.Vars(r)["genus"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data, err := json.Marshal(StrainJSON{Strain: strain})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Write(data)
}

func dbGetStrains(opt *StrainListOptions) ([]*Strain, error) {
	if opt == nil {
		return nil, errors.New("must provide options")
	}
	var vals []interface{}

	sql := `SELECT st.*, sp.species_name, array_agg(m.id) AS measurements,
		COUNT(m) AS total_measurements
		FROM strains st
		INNER JOIN species sp ON sp.id=st.species_id
		INNER JOIN genera g ON g.id=sp.genus_id AND LOWER(g.genus_name)=$1
		LEFT OUTER JOIN measurements m ON m.strain_id=st.id`
	vals = append(vals, opt.Genus)

	if len(opt.Ids) != 0 {
		var conds []string
		s := "st.id IN ("
		for i, id := range opt.Ids {
			s = s + fmt.Sprintf("$%v,", i+2) // start param index at 2
			vals = append(vals, id)
		}
		s = s[:len(s)-1] + ")"
		conds = append(conds, s)
		sql += " WHERE (" + strings.Join(conds, ") AND (") + ")"
	}

	sql += " GROUP BY st.id, sp.species_name;"

	var strains []*Strain
	err := DBH.Select(&strains, sql, vals...)
	if err != nil {
		return nil, err
	}
	return strains, nil
}

func dbGetStrain(id int64, genus string) (*Strain, error) {
	var strain Strain
	q := `SELECT st.*, sp.species_name, array_agg(m.id) AS measurements,
		COUNT(m) AS total_measurements
		FROM strains st
		INNER JOIN species sp ON sp.id=st.species_id
		INNER JOIN genera g ON g.id=sp.genus_id AND LOWER(g.genus_name)=$1
		LEFT OUTER JOIN measurements m ON m.strain_id=st.id
		WHERE st.id=$2
		GROUP BY st.id, sp.species_name;`
	if err := DBH.SelectOne(&strain, q, genus, id); err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrStrainNotFound
		}
		return nil, err
	}
	return &strain, nil
}

func dbUpdateStrain(strain *Strain) error {
	var species_id struct{ Id int64 }
	q := `SELECT id FROM species WHERE species_name = $1;`
	if err := DBH.SelectOne(&species_id, q, strain.SpeciesName); err != nil {
		return err
	}
	strain.StrainBase.SpeciesId = species_id.Id
	count, err := DBH.Update(strain.StrainBase)
	if err != nil {
		return err
	}
	if count != 1 {
		return ErrStrainNotUpdated
	}
	return nil
}
