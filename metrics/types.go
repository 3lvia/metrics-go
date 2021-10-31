package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

// Counter is a Metric that represents a single numerical value that only ever
// goes up.
type Counter interface {
	// Inc increments the counter by 1. Use Add to increment it by arbitrary
	// non-negative values.
	Inc()

	// Add adds the given value to the counter. It panics if the value is <
	// 0.
	Add(float64)
}

type counter struct {
	name    string
	labels  map[string]string
	inner   prometheus.Counter
	changes chan<- CountChange
}

func (c *counter) Inc() {
	c.inner.Inc()
	c.signal(1)
}

func (c *counter) Add(f float64) {
	c.inner.Add(f)
	c.signal(f)
}

func (c *counter) signal(v float64) {
	if c.inner == nil {
		return
	}
	if c.changes != nil {
		c.changes <- CountChange{
			Name:      c.name,
			Increment: v,
			Labels:    c.labels,
		}
	}
}

// Gauge is a Metric that represents a single numerical value that can
// arbitrarily go up and down.
type Gauge interface {
	// Set sets the Gauge to an arbitrary value.
	Set(float64)
	// Inc increments the Gauge by 1. Use Add to increment it by arbitrary
	// values.
	Inc()
	// Dec decrements the Gauge by 1. Use Sub to decrement it by arbitrary
	// values.
	Dec()
	// Add adds the given value to the Gauge. (The value can be negative,
	// resulting in a decrease of the Gauge.)
	Add(float64)
	// Sub subtracts the given value from the Gauge. (The value can be
	// negative, resulting in an increase of the Gauge.)
	Sub(float64)
}

type gauge struct {
	name    string
	labels  map[string]string
	inner   prometheus.Gauge
	changes chan<- GaugeChange
}

func (g *gauge) Set(f float64) {
	g.inner.Set(f)
	g.signal(f)
}

func (g *gauge) Inc() {
	g.inner.Inc()
	g.signal(1)
}

func (g *gauge) Dec() {
	g.inner.Dec()
	g.signal(-1)
}

func (g *gauge) Add(f float64) {
	g.inner.Add(f)
	g.signal(f)
}

func (g *gauge) Sub(f float64) {
	g.inner.Sub(f)
	g.signal(-1 * f)
}

func (g *gauge) signal(f float64) {
	if g.changes == nil {
		return
	}
	g.changes <- GaugeChange{
		Name:   g.name,
		Value:  f,
		Labels: g.labels,
	}
}

// CountChange represents a change to counter. This instance contains the increment, not the resulting value.
type CountChange struct {
	Name      string
	Increment float64
	Labels    map[string]string
}

// GaugeChange represents a change to a gauge. This instance contains the value that the gauge was changed by, not the
// resulting value.
type GaugeChange struct {
	Name   string
	Value  float64
	Labels map[string]string
}

