package temperature

import (
	"fmt"
	"github.com/reef-pi/reef-pi/controller/utils"
	"log"
	"time"
)

func (c *Controller) Check(tc TC) {
	if !tc.Enable {
		return
	}
	u := Usage{
		Time: utils.TeleTime(time.Now()),
	}
	reading, err := c.Read(tc)
	if err != nil {
		log.Println("ERROR: temperature sub-system. Failed to read  sensor. Error:", err)
		c.c.LogError("tc-"+tc.ID, "temperature sub-system. Failed to read  sensor "+tc.Name+". Error:"+err.Error())
		return
	}
	u.Temperature = reading
	log.Println("temperature sub-system:  sensor", tc.Name, "value:", reading)
	if tc.Control {
		if err := c.control(tc, &u); err != nil {
			log.Println("ERROR: Failed to execute temperature control logic. Error:", err)
		}
	}
	c.NotifyIfNeeded(tc, reading)
	c.statsMgr.Update(tc.ID, u)
}

func (c *Controller) control(tc TC, u *Usage) error {
	switch {
	case u.Temperature > tc.Max:
		log.Println("temperature subsystem: Current temperature is above maximum threshold. Executing cool down routine")
		u.Cooler += int(tc.Period)
		return c.coolDown(tc)
	case u.Temperature < tc.Min:
		log.Println("temparature subsystem: Current temperature is below minimum threshold. Executing warm up routine")
		u.Heater += int(tc.Period)
		return c.warmUp(tc)
	default:
		log.Println("temperature subsystem: Current temperature is within range switching off heater/cooler")
		c.switchOffAll(tc)
	}
	return nil
}

func (c *Controller) warmUp(tc TC) error {
	if tc.Cooler != "" {
		if err := c.equipment.Control(tc.Cooler, false); err != nil {
			return err
		}
	}
	if tc.Heater != "" {
		if err := c.equipment.Control(tc.Heater, true); err != nil {
			return err
		}
	}
	return nil
}

func (c *Controller) coolDown(tc TC) error {
	if tc.Heater != "" {
		if err := c.equipment.Control(tc.Heater, false); err != nil {
			return err
		}
	}
	if tc.Cooler != "" {
		if err := c.equipment.Control(tc.Cooler, true); err != nil {
			return err
		}
	}
	return nil
}

func (c *Controller) switchOffAll(tc TC) {
	if tc.Heater != "" {
		c.equipment.Control(tc.Heater, false)
	}
	if tc.Cooler != "" {
		c.equipment.Control(tc.Cooler, false)
	}
}

func (c *Controller) NotifyIfNeeded(tc TC, reading float64) {
	if !tc.Notify.Enable {
		return
	}
	subject := fmt.Sprintf("[Reef-Pi ALERT] temperature of '%s' out of range", tc.Name)
	format := "Current temperature (%f) is out of acceptable range ( %f -%f )"
	body := fmt.Sprintf(format, reading, tc.Notify.Min, tc.Notify.Max)
	if reading >= tc.Notify.Max {
		c.c.Telemetry().Alert(subject, "Tank is running hot."+body)
		return
	}
	if reading <= tc.Notify.Min {
		c.c.Telemetry().Alert(subject, "Tank is running cold. "+body)
		return
	}
}
