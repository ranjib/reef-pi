package connectors

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/reef-pi/reef-pi/controller/drivers"
	"github.com/reef-pi/reef-pi/controller/storage"
	"github.com/reef-pi/reef-pi/controller/utils"
)

func TestOutletsAPI(t *testing.T) {
	store, err := storage.TestDB()
	if err != nil {
		t.Fatal(err)
	}
	drvrs := drivers.TestDrivers(store)

	tr := utils.NewTestRouter()
	o := Outlet{Name: "Foo", Pin: 21, Driver: "rpi"}
	outlets := NewOutlets(drvrs, store)
	outlets.DevMode = true
	if err := outlets.Setup(); err != nil {
		t.Fatal(err)
	}
	outlets.LoadAPI(tr.Router)
	body := new(bytes.Buffer)
	json.NewEncoder(body).Encode(o)
	if err := tr.Do("PUT", "/api/outlets", body, nil); err != nil {
		t.Error(err)
	}
	body.Reset()
	json.NewEncoder(body).Encode(o)
	if err := tr.Do("POST", "/api/outlets/1", body, nil); err != nil {
		t.Error(err)
	}
	if err := tr.Do("GET", "/api/outlets/1", body, nil); err != nil {
		t.Error(err)
	}
	if err := tr.Do("GET", "/api/outlets", new(bytes.Buffer), nil); err != nil {
		t.Error(err)
	}
	if err := outlets.Configure("1", false); err != nil {
		t.Error(err)
	}
	o.Equipment = "1"
	if err := outlets.Update("1", o); err != nil {
		t.Error(err)
	}
	if err := outlets.Delete("1"); err == nil {
		t.Error("Expected to fail outlet deletion since equipment is attached to it")
	}
	body.Reset()
	o.Name = ""
	json.NewEncoder(body).Encode(o)
	if err := tr.Do("POST", "/api/outlets/1", body, nil); err == nil {
		t.Error("Expected to fail to update outlet since name is not set")
	}
	o.Equipment = ""
	o.Name = "asd"
	if err := outlets.Update("1", o); err != nil {
		t.Error(err)
	}
	if err := tr.Do("DELETE", "/api/outlets/1", new(bytes.Buffer), nil); err != nil {
		t.Error(err)
	}

	o.Name = ""
	if err := o.IsValid(drvrs); err == nil {
		t.Errorf("Outlet validation should fail if name is not set")
	}
	o.Name = "zsda"
	o.Pin = 1
	if err := o.IsValid(drvrs); err == nil {
		t.Errorf("Outlet validation should fail if GPIO pin is not valid")
	}
}
