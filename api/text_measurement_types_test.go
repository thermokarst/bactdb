package api

import (
	"testing"

	"github.com/thermokarst/bactdb/models"
)

func newTextMeasurementType() *models.TextMeasurementType {
	text_measurement_type := models.NewTextMeasurementType()
	return text_measurement_type
}

func TestTextMeasurementType_Get(t *testing.T) {
	setup()

	want := newTextMeasurementType()

	calledGet := false

	store.TextMeasurementTypes.(*models.MockTextMeasurementTypesService).Get_ = func(id int64) (*models.TextMeasurementType, error) {
		if id != want.Id {
			t.Errorf("wanted request for text_measurement_type %d but got %d", want.Id, id)
		}
		calledGet = true
		return want, nil
	}

	got, err := apiClient.TextMeasurementTypes.Get(want.Id)
	if err != nil {
		t.Fatal(err)
	}

	if !calledGet {
		t.Error("!calledGet")
	}
	if !normalizeDeepEqual(want, got) {
		t.Errorf("got %+v but wanted %+v", got, want)
	}
}

func TestTextMeasurementType_Create(t *testing.T) {
	setup()

	want := newTextMeasurementType()

	calledPost := false
	store.TextMeasurementTypes.(*models.MockTextMeasurementTypesService).Create_ = func(text_measurement_type *models.TextMeasurementType) (bool, error) {
		if !normalizeDeepEqual(want, text_measurement_type) {
			t.Errorf("wanted request for text_measurement_type %d but got %d", want, text_measurement_type)
		}
		calledPost = true
		return true, nil
	}

	success, err := apiClient.TextMeasurementTypes.Create(want)
	if err != nil {
		t.Fatal(err)
	}

	if !calledPost {
		t.Error("!calledPost")
	}
	if !success {
		t.Error("!success")
	}
}

func TestTextMeasurementType_List(t *testing.T) {
	setup()

	want := []*models.TextMeasurementType{newTextMeasurementType()}
	wantOpt := &models.TextMeasurementTypeListOptions{ListOptions: models.ListOptions{Page: 1, PerPage: 10}}

	calledList := false
	store.TextMeasurementTypes.(*models.MockTextMeasurementTypesService).List_ = func(opt *models.TextMeasurementTypeListOptions) ([]*models.TextMeasurementType, error) {
		if !normalizeDeepEqual(wantOpt, opt) {
			t.Errorf("wanted options %d but got %d", wantOpt, opt)
		}
		calledList = true
		return want, nil
	}

	text_measurement_types, err := apiClient.TextMeasurementTypes.List(wantOpt)
	if err != nil {
		t.Fatal(err)
	}

	if !calledList {
		t.Error("!calledList")
	}

	if !normalizeDeepEqual(&want, &text_measurement_types) {
		t.Errorf("got text_measurement_types %+v but wanted text_measurement_types %+v", text_measurement_types, want)
	}
}

func TestTextMeasurementType_Update(t *testing.T) {
	setup()

	want := newTextMeasurementType()

	calledPut := false
	store.TextMeasurementTypes.(*models.MockTextMeasurementTypesService).Update_ = func(id int64, text_measurement_type *models.TextMeasurementType) (bool, error) {
		if id != want.Id {
			t.Errorf("wanted request for text_measurement_type %d but got %d", want.Id, id)
		}
		if !normalizeDeepEqual(want, text_measurement_type) {
			t.Errorf("wanted request for text_measurement_type %d but got %d", want, text_measurement_type)
		}
		calledPut = true
		return true, nil
	}

	success, err := apiClient.TextMeasurementTypes.Update(want.Id, want)
	if err != nil {
		t.Fatal(err)
	}

	if !calledPut {
		t.Error("!calledPut")
	}
	if !success {
		t.Error("!success")
	}
}
