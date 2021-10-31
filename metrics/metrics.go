// Package metrics wraps prometheus.Metrics making it a bit easier to manage counters and gauges and also providing
// a mechanism by which all changes are communicated out via channels.
package metrics

import (
"fmt"
"github.com/prometheus/client_golang/prometheus"
"github.com/prometheus/client_golang/prometheus/promauto"
"sync"
)

// New creates and returns a new instance of Metrics.
func New(opts ...Option) Metrics {
	collector := &optionsCollector{}
	for _, opt := range opts {
		opt(collector)
	}
	return &metricsGuard{
		counterMap: map[string]Counter{},
		gaugeMap:   map[string]Gauge{},
		mux:        &sync.Mutex{},
		countChan:  collector.countChan,
		gaugeChan:  collector.gaugeChan,
	}
}

// Metrics provides the concrete abstraction used by clients of this package.
type Metrics interface {
	// Counter returns a counter for the given name and labels. Is useful in situation where the counter should be
	// increases by some other value that 1.
	Counter(name string, constLabels map[string]string) Counter

	// IncCounter increases the counter with the given name and labels by 1.
	IncCounter(name string, constLabels map[string]string)

	// Gauge returns a gauge with the given name.
	Gauge(name string) Gauge

	// SetGauge sets the given value in the gauge with the given name.
	SetGauge(name string, v int)
}

type metricsGuard struct {
	counterMap map[string]Counter
	gaugeMap   map[string]Gauge
	mux        *sync.Mutex
	countChan  chan<- CountChange
	gaugeChan  chan<- GaugeChange
}

func (g *metricsGuard) SetGauge(name string, v int) {
	gg := g.Gauge(name)
	gg.Set(float64(v))
}

func (g *metricsGuard) Gauge(name string) Gauge {
	if gg, ok := g.gaugeMap[name]; ok {
		return gg
	}

	g.mux.Lock()
	defer g.mux.Unlock()

	if gg, ok := g.gaugeMap[name]; ok {
		return gg
	}

	gg := &gauge{
		name:    name,
		labels:  nil,
		changes: g.gaugeChan,
		inner: promauto.NewGauge(prometheus.GaugeOpts{
			Name: name,
		}),
	}

	g.gaugeMap[name] = gg

	return gg
}

func (g *metricsGuard) IncCounter(name string, constLabels map[string]string) {
	c := g.Counter(name, constLabels)
	c.Inc()
}

func (g *metricsGuard) Counter(name string, constLabels map[string]string) Counter {
	key := metricsKey(name, constLabels)
	if c, ok := g.counterMap[key]; ok {
		return c
	}

	g.mux.Lock()
	defer g.mux.Unlock()

	if c, ok := g.counterMap[key]; ok {
		return c
	}

	c := &counter{
		name:    name,
		labels:  constLabels,
		changes: g.countChan,
		inner: promauto.NewCounter(prometheus.CounterOpts{
			Name:        name,
			ConstLabels: constLabels,
		}),
	}

	g.counterMap[key] = c

	return c
}

func metricsKey(name string, constLabels map[string]string) string {
	if constLabels == nil {
		return name
	}
	key := name
	for k, v := range constLabels {
		key = fmt.Sprintf("%s%s%s", key, k, v)
	}
	return key
}

