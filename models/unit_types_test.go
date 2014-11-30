package models

import (
	"net/http"
	"reflect"
	"testing"

	"github.com/thermokarst/bactdb/router"
)

func newUnitType() *UnitType {
	unit_type := NewUnitType()
	unit_type.Id = 1
	return unit_type
}

func TestUnitTypeService_Get(t *testing.T) {
	setup()
	defer teardown()

	want := newUnitType()

	var called bool
	mux.HandleFunc(urlPath(t, router.UnitType, map[string]string{"Id": "1"}), func(w http.ResponseWriter, r *http.Request) {
		called = true
		testMethod(t, r, "GET")

		writeJSON(w, want)
	})

	unit_type, err := client.UnitTypes.Get(want.Id)
	if err != nil {
		t.Errorf("UnitTypes.Get returned error: %v", err)
	}

	if !called {
		t.Fatal("!called")
	}

	normalizeTime(&want.CreatedAt, &want.UpdatedAt, &want.DeletedAt)

	if !reflect.DeepEqual(unit_type, want) {
		t.Errorf("UnitTypes.Get return %+v, want %+v", unit_type, want)
	}
}
