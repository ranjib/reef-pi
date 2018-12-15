package lighting

import (
	"encoding/json"
	"log"
	"math"
	"time"
)

type Profile struct {
	Type   string          `json:"type"`
	Config json.RawMessage `json:"config"`
}

type AutoConfig struct {
	Values []int `json:"values"` // 12 ticks after every 2 hours
}

type FixedConfig struct {
	Value int `json:"value"`
}

type DiurnalConfig struct {
	Start string `json:"start"`
	End   string `json:"end"`
}

const TimeFormat = "15:04"

func (ch Channel) GetValueDiurnal(t time.Time) float64 {
	var d DiurnalConfig
	if err := json.Unmarshal(ch.Profile.Config, &d); err != nil {
		log.Println("ERROR: lighting subsysten failed to typecast diurnal config. Error:", err)
		return 0
	}

	s, err := time.Parse(TimeFormat, d.Start)
	if err != nil {
		log.Println("ERROR: lighting subsystem, failed to parse start time in diurnal cycle. Error:", err)
		return 0
	}
	e, err := time.Parse(TimeFormat, d.End)
	if err != nil {
		log.Println("ERROR: lighting subsystem, failed to parse end time in diurnal cycle. Error:", err)
		return 0
	}
	s = time.Date(t.Year(), t.Month(), t.Day(), s.Hour(), s.Minute(), s.Second(), 0, t.Location())
	e = time.Date(t.Year(), t.Month(), t.Day(), e.Hour(), e.Minute(), e.Second(), 0, t.Location())
	if e.Before(s) {
		e = e.Add(time.Hour * 24)
		if t.Before(s) {
			t = t.Add(time.Hour * 24)
		}
	}
	//log.Println("start:", s, "end:", e, "t:", t)
	if t.Before(s) {
		log.Println("Lighting sub-system: diurnal profile: before start time. Value: 0")
		return 0
	}
	if t.After(e) {
		log.Println("Lighting sub-system: diurnal profile: after end time. Value: 0")
		return 0
	}
	totalMinutes := int(e.Sub(s) / time.Minute)
	pastMinutes := int(t.Sub(s) / time.Minute)
	percent := float64(pastMinutes) * 2 * math.Pi / float64(totalMinutes)
	k := math.Pow(math.Cos(percent), 3)
	v := (1 - k) * float64(ch.Max-ch.Min)
	v = v + float64(ch.Min)
	if v > float64(ch.Max) {
		v = float64(ch.Max)
	}
	return v
}

type FreeFormConfig struct {
	Values []int `json:"values"`
	Period int   `json:"period"`
}

func (ch Channel) GetValue(t time.Time) float64 {
	switch ch.Profile.Type {
	case "diurnal":
		return ch.GetValueDiurnal(t)
	case "auto":
		return ch.GetValueAuto(t)
	case "fixed":
		return ch.GetValueFixed()
	default:
		return 0
	}
}

func (ch Channel) GetValueFixed() float64 {
	var f FixedConfig
	if err := json.Unmarshal(ch.Profile.Config, &f); err != nil {
		log.Println("ERROR: lighting subsysten failed to typecast fixed config. Error", err)
		return 0
	}
	return float64(f.Value)
}
func (ch Channel) GetValueAuto(t time.Time) float64 {
	var a AutoConfig
	if err := json.Unmarshal(ch.Profile.Config, &a); err != nil {
		log.Println("ERROR: lighting subsysten failed to typecast auto config. Error", err)
		return 0
	}
	series := a.Values
	l := len(series)
	if l < 12 {
		for i := l - 1; i < 12; i++ {
			series = append(series, 0)
		}
	}
	h1 := t.Hour() / 2
	h2 := h1 + 1
	if h2 >= 12 {
		h2 = 0
	}
	m := float64(t.Minute() + ((t.Hour() % 2) * 60))
	from := float64(series[h1])
	to := float64(series[h2])
	return from + ((to - from) / 120.0 * m)
}
