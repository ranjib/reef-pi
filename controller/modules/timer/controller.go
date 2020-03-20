package timer

import (
	"encoding/json"
	"fmt"
	"log"

	cron "github.com/robfig/cron/v3"

	"github.com/reef-pi/reef-pi/controller"
	"github.com/reef-pi/reef-pi/controller/storage"
)

const _cronParserSpec = cron.Second | cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow | cron.Descriptor

type Controller struct {
	runner    *cron.Cron
	cronIDs   map[string]cron.EntryID
	equipment controller.Subsystem
	c         controller.Controller
	macro     controller.Subsystem
}

func New(c controller.Controller) *Controller {
	return &Controller{
		cronIDs: make(map[string]cron.EntryID),
		runner:  cron.New(cron.WithParser(cron.NewParser(_cronParserSpec))),
		c:       c,
	}
}

func (c *Controller) Setup() error {
	if macro, err := c.c.Subsystem(storage.MacroBucket); err == nil {
		c.macro = macro
	}
	if eq, err := c.c.Subsystem(storage.EquipmentBucket); err == nil {
		c.equipment = eq
	}
	return c.c.Store().CreateBucket(Bucket)
}

func (c *Controller) Start() {
	if err := c.loadAllJobs(); err != nil {
		log.Println("ERROR: timer-subsystem: Failed to load timer jobs. Error:", err)
	}
	c.runner.Start()
}

func (c *Controller) Stop() {
	c.runner.Stop()
}

func (c *Controller) On(id string, on bool) error {
	j, err := c.Get(id)
	if err != nil {
		return err
	}
	j.Enable = on
	return c.Update(id, j)
}

func (c *Controller) InUse(depType, id string) ([]string, error) {
	var deps []string
	switch depType {
	case storage.EquipmentBucket:
		ts, err := c.List()
		if err != nil {
			return deps, err
		}
		for _, timer := range ts {
			if timer.Type == depType {
				var ue UpdateEquipment
				if err := json.Unmarshal(timer.Target, &ue); err != nil {
					return deps, err
				}
				if ue.ID == id {
					deps = append(deps, timer.Name)
				}
			}
		}
		return deps, nil
	case storage.MacroBucket:
		ts, err := c.List()
		if err != nil {
			return deps, err
		}
		for _, timer := range ts {
			if timer.Type == depType {
				var m TriggerMacro
				if err := json.Unmarshal(timer.Target, &m); err != nil {
					return deps, err
				}
				if m.ID == id {
					deps = append(deps, timer.Name)
				}
			}
		}
		return deps, nil
	default:
		return deps, fmt.Errorf("unknown dep type:%s", depType)
	}
}
