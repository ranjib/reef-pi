package controller

import (
	"fmt"
	"github.com/hybridgroup/gobot/platforms/raspi"
)

type RaspiConfig struct {
	Relay1 string      `yaml:"relay_1"`
	Relay2 string      `yaml:"relay_2"`
	Doser1 DoserConfig `yaml:"doser_1"`
	Doser2 DoserConfig `yaml:"doser_2"`
}

func DefaultRaspiConfig() RaspiConfig {
	var config RaspiConfig
	config.Relay1 = "1"
	config.Relay2 = "2"
	config.Doser1.PWMPin = "3"
	config.Doser1.IN1Pin = "4"
	config.Doser1.IN2Pin = "5"
	config.Doser2.PWMPin = "6"
	config.Doser2.IN1Pin = "7"
	config.Doser2.IN2Pin = "8"
	return config
}

type Raspi struct {
	config  *RaspiConfig
	devices map[string]Device
}

func (c *Raspi) GetDevice(name string) (Device, error) {
	dev, ok := c.devices[name]
	if !ok {
		return nil, fmt.Errorf("No such device: '%s'", name)
	}
	return dev, nil
}

func NewRaspi(config *RaspiConfig) *Raspi {
	return &Raspi{
		devices: loadDevices(config),
	}
}

func loadDevices(config *RaspiConfig) map[string]Device {
	conn := raspi.NewRaspiAdaptor("raspi")
	devices := make(map[string]Device)
	devices["relay_1"] = NewRelay("relay_1", conn, config.Relay1)
	devices["relay_2"] = NewRelay("relay_2", conn, config.Relay2)
	devices["doser_1"] = NewDoser("doser_1", conn, config.Doser1)
	devices["doser_2"] = NewDoser("doser_2", conn, config.Doser2)
	return devices
}
