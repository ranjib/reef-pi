package connectors

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/reef-pi/reef-pi/controller/drivers"

	"github.com/reef-pi/reef-pi/controller/storage"
	"github.com/reef-pi/reef-pi/controller/utils"
)

func TestAnalogInputsAPI(t *testing.T) {
	store, err := storage.TestDB()
	if err != nil {
		t.Fatal(err)
	}
	tr := utils.NewTestRouter()

	if err != nil {
		t.Error(err)
	}
	drvrs := drivers.TestDrivers(store)
	d1 := drivers.Driver{
		Name:   "lighting",
		Type:   "pca9685",
		Config: []byte(`{"address":64, "frequency":1000}`),
	}
	if err := drvrs.Create(d1); err != nil {
		t.Fatal(err)
	}
	j := AnalogInput{Name: "Foo", Pin: 0, Driver: "rpi"}
	ais := NewAnalogInputs(drvrs, store)
	if err := ais.Setup(); err != nil {
		t.Fatal(err)
	}
	ais.LoadAPI(tr.Router)
	body := new(bytes.Buffer)
	json.NewEncoder(body).Encode(j)
	if err := tr.Do("PUT", "/api/analog_inputs", body, nil); err != nil {
		t.Error(err)
	}

	body.Reset()
	j.Driver = "1"
	json.NewEncoder(body).Encode(j)
	if err := tr.Do("POST", "/api/analog_inputs/1", body, nil); err != nil {
		t.Error(err)
	}

	body.Reset()
	j.Name = ""
	json.NewEncoder(body).Encode(j)
	if err := tr.Do("PUT", "/api/analog_inputs", body, nil); err == nil {
		t.Error("AnalogInput creation expected to fail when analog_input name is absent")
	}
	body.Reset()
	j.Name = "zd"
	json.NewEncoder(body).Encode(j)
	if err := tr.Do("PUT", "/api/analog_inputs", body, nil); err == nil {
		t.Error("AnalogInput creation expected to fail when analog_input pin is empty")
	}
	body.Reset()
	j.Pin = 16
	json.NewEncoder(body).Encode(j)
	if err := tr.Do("PUT", "/api/analog_inputs", body, nil); err == nil {
		t.Error("AnalogInput creation expected to fail when pca9685 pin is invalid (not 0-14)")
	}
	body.Reset()
	j.Driver = "rpi"
	j.Pin = 3
	json.NewEncoder(body).Encode(j)
	if err := tr.Do("PUT", "/api/analog_inputs", body, nil); err == nil {
		t.Error("AnalogInput creation expected to fail when rpi pin is invalid (not 0 or 1)")
	}
	body.Reset()
	j.Driver = ""
	json.NewEncoder(body).Encode(j)
	if err := tr.Do("POST", "/api/analog_inputs/1", body, nil); err == nil {
		t.Error("AnalogInput updateexpected to fail when driver is invalid (rpi and pca9685 are only valid values)")
	}

	if err := tr.Do("GET", "/api/analog_inputs/1", body, nil); err != nil {
		t.Error(err)
	}
	if err := tr.Do("GET", "/api/analog_inputs", new(bytes.Buffer), nil); err != nil {
		t.Error(err)
	}
	x := AnalogReading{Value: 12.23}
	body.Reset()
	json.NewEncoder(body).Encode(&x)
	if err := tr.Do("POST", "/api/analog_inputs/1/read", body, nil); err != nil {
		t.Error(err)
	}
	if err := tr.Do("DELETE", "/api/analog_inputs/1", new(bytes.Buffer), nil); err != nil {
		t.Error(err)
	}
}
