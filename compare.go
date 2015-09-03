package main

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"
)

func handleCompare(w http.ResponseWriter, r *http.Request) *appError {
	// types
	type Comparisions map[string]map[string]string
	type ComparisionsJSON [][]string

	// vars
	mimeType := r.FormValue("mimeType")
	if mimeType == "" {
		mimeType = "json"
	}
	claims := getClaims(r)
	var header string
	var data []byte

	// Get measurements for comparision
	measService := MeasurementService{}
	opt := r.URL.Query()
	opt.Del("mimeType")
	opt.Del("token")
	opt.Add("Genus", mux.Vars(r)["genus"])
	measurementsEntity, appErr := measService.list(&opt, &claims)
	if appErr != nil {
		return appErr
	}
	measurementsPayload := (measurementsEntity).(*MeasurementsPayload)

	// Assemble matrix
	characteristic_ids := strings.Split(opt.Get("characteristic_ids"), ",")
	strain_ids := strings.Split(opt.Get("strain_ids"), ",")

	comparisions := make(Comparisions)
	for _, characteristic_id := range characteristic_ids {
		characteristic_id_int, _ := strconv.ParseInt(characteristic_id, 10, 0)
		values := make(map[string]string)
		for _, strain_id := range strain_ids {
			strain_id_int, _ := strconv.ParseInt(strain_id, 10, 0)
			for _, m := range *measurementsPayload.Measurements {
				if (m.CharacteristicId == characteristic_id_int) && (m.StrainId == strain_id_int) {
					values[strain_id] = m.Value()
				}
			}
		}

		comparisions[characteristic_id] = values
	}

	// Return, based on mimetype
	switch mimeType {
	case "json":
		header = "application/json"

		comparisionsJSON := make(ComparisionsJSON, 0)
		for _, characteristic_id := range characteristic_ids {
			row := []string{characteristic_id}
			for _, strain_id := range strain_ids {
				row = append(row, comparisions[characteristic_id][strain_id])
			}
			comparisionsJSON = append(comparisionsJSON, row)
		}

		data, _ = json.Marshal(comparisionsJSON)
	case "csv":
		header = "text/csv"

		// maps to translate ids
		strains := make(map[string]string)
		for _, strain := range *measurementsPayload.Strains {
			strains[fmt.Sprintf("%d", strain.Id)] = fmt.Sprintf("%s (%s)", strain.SpeciesName(), strain.StrainName)
		}
		characteristics := make(map[string]string)
		for _, characteristic := range *measurementsPayload.Characteristics {
			characteristics[fmt.Sprintf("%d", characteristic.Id)] = characteristic.CharacteristicName
		}

		b := &bytes.Buffer{}
		wr := csv.NewWriter(b)

		// Write header row
		r := []string{"Characteristic"}
		for _, strain_id := range strain_ids {
			r = append(r, strains[strain_id])
		}
		wr.Write(r)

		// Write data
		for key, record := range comparisions {
			r := []string{characteristics[key]}
			for _, val := range record {
				r = append(r, val)
			}
			wr.Write(r)
		}
		wr.Flush()

		data = b.Bytes()

		w.Header().Set("Content-Disposition", fmt.Sprintf(`inline; filename="compare-%d.csv"`, int32(time.Now().Unix())))
	}

	// Wrap it up
	w.Header().Set("Content-Type", header)
	w.Write(data)
	return nil
}
