package metrics

type optionsCollector struct {
	countChan chan<- CountChange
	gaugeChan chan<- GaugeChange
}

// Option for configuring this package.
type Option func(collector *optionsCollector)

// WithOutputChannels sets channels which will trigger whenever a counter of gauge is altered. This feature is meant
// for testing purposes.
func WithOutputChannels(countChan chan<- CountChange, gaugeChan chan<- GaugeChange) Option {
	return func(collector *optionsCollector) {
		collector.countChan = countChan
		collector.gaugeChan = gaugeChan
	}
}