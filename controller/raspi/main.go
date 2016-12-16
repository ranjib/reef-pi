package raspi

import (
	"fmt"
	"github.com/boltdb/bolt"
	pi "github.com/hybridgroup/gobot/platforms/raspi"
	"github.com/ranjib/reefer/controller"
	"log"
	"time"
)

type Raspi struct {
	db           *bolt.DB
	conn         *pi.RaspiAdaptor
	schedules    map[controller.Device]controller.Scheduler
	modules      map[string]controller.Module
	lighting     *Lighting
	deviceAPI    controller.CrudAPI
	boardAPI     controller.CrudAPI
	outletAPI    controller.CrudAPI
	jobAPI       *JobAPI
	equipmentAPI controller.CrudAPI
}

func (r *Raspi) Name() string {
	return "raspberry-pi"
}

func (c *Raspi) GetModule(name string) (controller.Module, error) {
	module, ok := c.modules[name]
	if !ok {
		return nil, fmt.Errorf("No such module: '%s'", name)
	}
	return module, nil
}

func New() (*Raspi, error) {
	db, err := bolt.Open("reefer.db", 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		return nil, err
	}
	conn := pi.NewRaspiAdaptor("raspi")
	deviceAPI, err := NewDeviceAPI(conn, db)
	if err != nil {
		return nil, err
	}
	boardAPI, err := NewBoardAPI(db)
	if err != nil {
		return nil, err
	}
	outletAPI, err := NewOutletAPI(db)
	if err != nil {
		return nil, err
	}
	jobAPI, err := NewJobAPI(conn, db)
	if err != nil {
		return nil, err
	}
	equipmentAPI, err := NewEquipmentAPI(conn, db)
	if err != nil {
		return nil, err
	}
	r := &Raspi{
		db:           db,
		conn:         conn,
		deviceAPI:    deviceAPI,
		outletAPI:    outletAPI,
		boardAPI:     boardAPI,
		jobAPI:       jobAPI,
		schedules:    make(map[controller.Device]controller.Scheduler),
		lighting:     NewLighting(),
		equipmentAPI: equipmentAPI,
	}
	return r, nil
}

func (r *Raspi) Schedule(dev controller.Device, sched controller.Scheduler) error {
	if _, ok := r.schedules[dev]; ok {
		return fmt.Errorf("Device %s already scheduled", dev.Name())
	}
	log.Printf("Added %s[ %s]\n", sched.Name(), dev.Name())
	r.schedules[dev] = sched
	if !sched.IsRunning() {
		go sched.Start(dev)
	}
	return nil
}

func (r *Raspi) Start() error {
	for dev, sched := range r.schedules {
		go sched.Start(dev)
	}
	r.jobAPI.Start()
	log.Println("Started Controller:", r.Name())
	return nil
}

func (r *Raspi) Stop() error {
	for _, sched := range r.schedules {
		sched.Stop()
	}
	defer r.jobAPI.Stop()
	defer r.db.Close()
	log.Println("Stopped Controller:", r.Name())
	return nil
}
