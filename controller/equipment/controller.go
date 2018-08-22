package equipment

import (
	"github.com/kidoman/embd"
	"github.com/reef-pi/reef-pi/controller/connectors"
	"github.com/reef-pi/reef-pi/controller/types"
	"log"
)

type Config struct {
	DevMode bool `json:"dev_mode"`
}
type Check func(string) (bool, error)
type Controller struct {
	config    Config
	telemetry types.Telemetry
	store     types.Store
	outlets   *connectors.Outlets
	checks    []Check
}

func New(config Config, outlets *connectors.Outlets, store types.Store, telemetry types.Telemetry) *Controller {
	return &Controller{
		config:    config,
		telemetry: telemetry,
		store:     store,
		outlets:   outlets,
		checks:    []Check{},
	}
}

func (c *Controller) Setup() error {
	return c.store.CreateBucket(Bucket)
}

func (c *Controller) Start() {
	if !c.config.DevMode {
		embd.InitGPIO()
	}
	eqs, err := c.List()
	if err != nil {
		log.Println("ERROR: equipment subsystem: failed to list equipment. Error:", err)
		return
	}
	for _, eq := range eqs {
		if err := c.outlets.Configure(eq.Outlet, eq.On); err != nil {
			log.Println("ERROR: equipment subsystem: Failed to sync equipment", eq.Name, ". Error:", err)
		}
	}
	log.Println("INFO: equipment subsystem: Finished syncing all equipment")
}

func (c *Controller) Stop() {
	if c.config.DevMode {
		log.Println("Equipment subsystem is running in dev mode, skipping GPIO closing")
		return
	}
	embd.CloseGPIO()
}

func (c *Controller) AddCheck(check Check) {
	c.checks = append(c.checks, check)
	return
}

func (c *Controller) IsEquipmentInUse(id string) (bool, error) {
	for i, checkFn := range c.checks {
		inUse, err := checkFn(id)
		if err != nil {
			log.Println("ERROR: equipment subsystem: Equipment in use check returned error. Error:", err)
			return true, err
		}
		if inUse {
			log.Println("DEBUG: equipment subsystem: Equipment in use returned true from:", i)
			return true, nil
		}
	}
	return false, nil
}

func (c *Controller) On(id string, b bool) error {
	e, err := c.Get(id)
	if err != nil {
		return err
	}
	e.On = b
	return c.Update(id, e)
}
