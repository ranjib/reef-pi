package utils

import (
	"encoding/json"
	"github.com/reef-pi/reef-pi/controller/types"
	"testing"
)

type testMetric struct {
	Value int
}

func (m1 *testMetric) Before(m types.Metric) bool {
	return true
}

func (m1 *testMetric) Rollup(m types.Metric) (types.Metric, bool) {
	return m1, true
}

func TestStatsManager(t *testing.T) {
	store, err := TestDB()
	if err != nil {
		t.Error(err)
	}
	store.CreateBucket("test-subsystem")
	mgr := NewStatsManager(store, "test-subsystem", 10, 10)
	metric := &testMetric{}
	mgr.Update("foo", metric)
	mgr.Update("foo", metric)
	if err := mgr.Save("foo"); err != nil {
		t.Error(err)
	}
	fn := func(_ json.RawMessage) interface{} { return metric }
	if err := mgr.Load("foo", fn); err != nil {
		t.Error(err)
	}
}
