package daemon

import (
	"github.com/reef-pi/reef-pi/controller/settings"
	"github.com/reef-pi/reef-pi/controller/storage"

	"net/http"
	"testing"
)

func TestReefPi(t *testing.T) {
	http.DefaultServeMux = new(http.ServeMux)
	conf, err := ParseConfig("../../build/config.yaml")
	if err != nil {
		t.Fatal("Failed to parse example config file. Error:", err)
	}
	conf.Database = "reef-pi.db"
	store, err := storage.NewStore(conf.Database)
	defer store.Close()

	if err != nil {
		t.Fatal(err)
	}
	initializeSettings(store)
	s := settings.DefaultSettings
	s.Capabilities.DevMode = true
	if err := store.Update(Bucket, "settings", s); err != nil {
		t.Fatal(err)
	}
	store.Close()
	r, err := New("0.1", conf.Database)
	if err != nil {
		t.Fatal("Failed to create new reef-pi controller. Error:", err)
	}
	r.settings.Capabilities.DevMode = true
	r.settings.Capabilities.Doser = true
	r.settings.Capabilities.Lighting = true
	r.settings.Capabilities.Camera = true
	r.settings.Capabilities.Equipment = true
	r.settings.Capabilities.Timers = true
	r.settings.Capabilities.Ph = true
	if err := r.Start(); err != nil {
		t.Fatal("Failed to load subsystem. Error:", err)
	}
	if err := r.API(); err != nil {
		t.Error(err)
	}
	if _, err := r.Subsystem("timers"); err != nil {
		t.Error(err)
	}
	if _, err := r.Subsystem("invalid"); err == nil {
		t.Errorf("invalid subsystem fetch should fail")
	}
	if err := r.Stop(); err != nil {
		t.Fatal(err)
	}

}
