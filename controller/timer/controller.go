package timer

import (
	"github.com/reef-pi/reef-pi/controller/equipment"
	"github.com/reef-pi/reef-pi/controller/utils"
	"gopkg.in/robfig/cron.v2"
	"log"
)

type Controller struct {
	store     utils.Store
	runner    *cron.Cron
	cronIDs   map[string]cron.EntryID
	telemetry *utils.Telemetry
	equipment *equipment.Controller
}

func New(store utils.Store, telemetry *utils.Telemetry, e *equipment.Controller) *Controller {
	return &Controller{
		cronIDs:   make(map[string]cron.EntryID),
		telemetry: telemetry,
		store:     store,
		equipment: e,
	}
}

func (c *Controller) IsEquipmentInUse(id string) (bool, error) {
	jobs, err := c.List()
	if err != nil {
		return false, err
	}
	for _, j := range jobs {
		if j.Equipment == id {
			return true, nil
		}
	}
	return false, nil
}

func (c *Controller) Setup() error {
	return c.store.CreateBucket(Bucket)
}

func (c *Controller) Start() {
	c.runner = cron.New()
	if err := c.loadAllJobs(); err != nil {
		log.Println("ERROR: timer-subsystem: Failed to load timer jobs. Error:", err)
	}
	c.runner.Start()
}

func (c *Controller) Stop() {
	c.runner.Stop()
	c.runner = nil
}
