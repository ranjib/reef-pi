package macro

import (
	"fmt"
	"github.com/reef-pi/reef-pi/controller/types"
	"github.com/reef-pi/reef-pi/controller/utils"
	"sync"
)

const Bucket = types.MacroBucket
const UsageBucket = types.MacroUsageBucket

type Subsystem struct {
	sync.Mutex
	telemetry  types.Telemetry
	store      types.Store
	devMode    bool
	quitters   map[string]chan struct{}
	statsMgr   *utils.StatsManager
	controller types.Controller
}

func New(devMode bool, c types.Controller, store types.Store, telemetry types.Telemetry) (*Subsystem, error) {
	return &Subsystem{
		telemetry:  telemetry,
		store:      store,
		devMode:    devMode,
		controller: c,
	}, nil
}

func (s *Subsystem) Setup() error {
	return s.store.CreateBucket(Bucket)
}

func (s *Subsystem) Start() {
}

func (s *Subsystem) Stop() {
}
func (s *Subsystem) On(id string, b bool) error {
	return fmt.Errorf("Macro sub system does not support 'on' API yet")
}
