package controller

import (
	"github.com/reef-pi/reef-pi/controller/types"
	"github.com/reef-pi/reef-pi/controller/utils"
	"log"
	"net/http"
	"os"
)

type Settings struct {
	Name            string            `json:"name"`
	Interface       string            `json:"interface"`
	Address         string            `json:"address"`
	Display         bool              `json:"display"`
	Notification    bool              `json:"notification"`
	Capabilities    Capabilities      `json:"capabilities"`
	HealthCheck     HealthCheckNotify `json:"health_check"`
	HTTPS           bool              `json:"https"`
	PCA9685         bool              `json:"pca9685"`
	Pprof           bool              `json:"pprof"`
	RPI_PWMFreq     int               `json:"rpi_pwm_freq"`
	PCA9685_PWMFreq int               `json:"pca9685_pwm_freq"`
}

var DefaultSettings = Settings{
	Name:            "reef-pi",
	Interface:       "wlan0",
	Address:         "localhost:80",
	Capabilities:    DefaultCapabilities,
	RPI_PWMFreq:     100,
	PCA9685_PWMFreq: 1500,
	HealthCheck: HealthCheckNotify{
		MaxMemory: 500,
		MaxCPU:    2,
	},
}

func loadSettings(store types.Store) (Settings, error) {
	var s Settings
	if err := store.Get(Bucket, "settings", &s); err != nil {
		return s, err
	}
	return s, nil
}

func initializeSettings(store types.Store) (Settings, error) {
	if os.Getenv("DEV_MODE") == "1" {
		DefaultSettings.Capabilities.DevMode = true
		DefaultSettings.Capabilities.Dashboard = true
		DefaultSettings.Capabilities.Equipment = true
		DefaultSettings.Capabilities.Timers = true
		DefaultSettings.Capabilities.Lighting = true
		DefaultSettings.Capabilities.Temperature = true
		DefaultSettings.Capabilities.ATO = true
		DefaultSettings.Capabilities.Macro = true
		DefaultSettings.Capabilities.Doser = true
		DefaultSettings.Capabilities.Ph = true

		DefaultSettings.Address = "0.0.0.0:8080"
		log.Println("DEV_MODE environment variable set. Turning on dev_mode. Address set to localhost:8080")
	}
	if err := store.CreateBucket(Bucket); err != nil {
		log.Println("ERROR:Failed to create bucket:", Bucket, ". Error:", err)
		return DefaultSettings, err
	}
	return DefaultSettings, store.Update(Bucket, "settings", DefaultSettings)
}

func (r *ReefPi) GetSettings(w http.ResponseWriter, req *http.Request) {
	fn := func(_ string) (interface{}, error) {
		var s Settings
		return &s, r.store.Get(Bucket, "settings", &s)
	}
	utils.JSONGetResponse(fn, w, req)
}

func (r *ReefPi) UpdateSettings(w http.ResponseWriter, req *http.Request) {
	var s Settings
	fn := func(_ string) error {
		return r.store.Update(Bucket, "settings", s)
	}
	utils.JSONUpdateResponse(&s, fn, w, req)
}
