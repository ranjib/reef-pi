package lighting

import (
	"fmt"
	"log"
	"time"
)

func ValidateValues(values []int) error {
	if len(values) != 12 {
		return fmt.Errorf("Expect 12 values instead of: %d", len(values))
	}
	for i, v := range values {
		if (v < 0) || (v > 100) {
			return fmt.Errorf(" value %d on index %d is out of range (0-99)", v, i)
		}
	}
	return nil
}

func (c *Controller) StartCycle() {
	ticker := time.NewTicker(c.config.Interval)
	log.Println("Starting lighting cycle")
	c.mu.Lock()
	c.running = true
	c.mu.Unlock()
	c.syncLights()
	for {
		select {
		case <-c.stopCh:
			ticker.Stop()
			return
		case <-ticker.C:
			c.syncLights()
		}
	}
}

func (c *Controller) syncLights() {
	lights, err := c.List()
	if err != nil {
		log.Println("ERROR: lighting sub-system:  Failed to list lights. Error:", err)
		return
	}
	for _, light := range lights {
		c.syncLight(light)
	}
}

func (c *Controller) StopCycle() {
	c.mu.Lock()
	if !c.running {
		log.Println("lighting subsystem: controller is not running.")
		return
	}
	c.mu.Unlock()
	c.stopCh <- struct{}{}
	c.running = false
	log.Println("Stopped lighting cycle")
}

func (c *Controller) UpdateChannel(jack string, ch Channel, v float64) {
	if ch.Reverse {
		v = 100 - v
	}
	log.Println("lighting-subsystem: Setting PWM value:", v, " at channel:", ch.Pin)
	pv := make(map[int]float64)
	pv[ch.Pin] = v
	if err := c.jacks.Control(jack, pv); err != nil {
		log.Println("ERROR: lighting-subsystem: Failed to set pwm value. Error:", err)
	}
}
