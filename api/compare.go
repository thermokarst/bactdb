package api

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/thermokarst/bactdb/Godeps/_workspace/src/github.com/gorilla/mux"
	"github.com/thermokarst/bactdb/helpers"
	"github.com/thermokarst/bactdb/payloads"
	"github.com/thermokarst/bactdb/types"
)

// HandleCompare is a HTTP handler for comparision.
// Comparision requires a list of strain ids and a list of characteristic ids.
// The id order dictates the presentation order.
func HandleCompare(w http.ResponseWriter, r *http.Request) *types.AppError {
	// types
	type Comparisions map[string]map[string]string
	type ComparisionsJSON [][]string

	// vars
	mimeType := r.FormValue("mimeType")
	if mimeType == "" {
		mimeType = "json"
	}
	claims := helpers.GetClaims(r)
	var header string
	var data []byte

	// Get measurements for comparision
	measService := MeasurementService{}
	opt := r.URL.Query()
	opt.Del("mimeType")
	opt.Del("token")
	opt.Add("Genus", mux.Vars(r)["genus"])
	measurementsEntity, appErr := measService.List(&opt, &claims)
	if appErr != nil {
		return appErr
	}
	measurementsPayload := (measurementsEntity).(*payloads.Measurements)

	// Assemble matrix
	characteristicIDs := strings.Split(opt.Get("characteristic_ids"), ",")
	strainIDs := strings.Split(opt.Get("strain_ids"), ",")

	comparisions := make(Comparisions)
	for _, characteristicID := range characteristicIDs {
		characteristicIDInt, _ := strconv.ParseInt(characteristicID, 10, 0)
		values := make(map[string]string)
		for _, strainID := range strainIDs {
			strainIDInt, _ := strconv.ParseInt(strainID, 10, 0)
			for _, m := range *measurementsPayload.Measurements {
				if (m.CharacteristicID == characteristicIDInt) && (m.StrainID == strainIDInt) {
					if m.Notes.Valid {
						values[strainID] = fmt.Sprintf("%s (%s)", m.Value(), m.Notes.String)
					} else {
						if values[strainID] != "" {
							values[strainID] = fmt.Sprintf("%s, %s", values[strainID], m.Value())
						} else {
							values[strainID] = m.Value()
						}
					}
				}
			}
			// If the strain doesn't have a measurement for this characteristic,
			// stick an empty value in anyway (for CSV).
			if _, ok := values[strainID]; !ok {
				values[strainID] = ""
			}
		}

		comparisions[characteristicID] = values
	}

	// Return, based on mimetype
	switch mimeType {
	case "json":
		header = "application/json"

		comparisionsJSON := make(ComparisionsJSON, 0)
		for _, characteristicID := range characteristicIDs {
			row := []string{characteristicID}
			for _, strainID := range strainIDs {
				row = append(row, comparisions[characteristicID][strainID])
			}
			comparisionsJSON = append(comparisionsJSON, row)
		}

		data, _ = json.Marshal(comparisionsJSON)
	case "csv":
		header = "text/csv"

		// maps to translate ids
		strains := make(map[string]string)
		for _, strain := range *measurementsPayload.Strains {
			strains[fmt.Sprintf("%d", strain.ID)] = fmt.Sprintf("%s (%s)", strain.SpeciesName(), strain.StrainName)
		}
		characteristics := make(map[string]string)
		for _, characteristic := range *measurementsPayload.Characteristics {
			characteristics[fmt.Sprintf("%d", characteristic.ID)] = characteristic.CharacteristicName
		}

		b := &bytes.Buffer{}
		wr := csv.NewWriter(b)

		// Write header row
		r := []string{"Characteristic"}
		for _, strainID := range strainIDs {
			r = append(r, strains[strainID])
		}
		wr.Write(r)

		// Write data
		for _, characteristicID := range characteristicIDs {
			r := []string{characteristics[characteristicID]}
			for _, strainID := range strainIDs {
				r = append(r, comparisions[characteristicID][strainID])
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
