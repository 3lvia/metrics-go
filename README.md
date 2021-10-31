# metrics-go
This package wraps prometheus.Metrics making it a bit easier to manage counters and gauges and also providing mechanism by which all changes are communicated out via channels.

## Standard usage
```
import "github.com/3lvia/metrics-go/metrics"

m := metrics.New()

m.IncCounter("my-counter", metrics.DayLabels())
```

## Usage with channels (mostly for testing purposes)
```
import (
    "fmt"
    "github.com/3lvia/metrics-go/metrics"
)

countChanges := make(chan metrics.CountChange)
gaugeChanges := make(chan metrics.GaugeChange)
m := metrics.New(metrics.WithOutputChannels(countChanges, gaugeChanges))

go func(cc <-chan metrics.CountChange, gc <-chan metrics.GaugeChan) {
    for {
        select {
            case c := <- cc:
                fmt.Printf("count change %s received", c.Name)
            case g := <- gc:
                fmt.Printf("gauge change %s received", g.Name)
        }
    }
} (countChanges, gaugeChanges)

m.IncCounter("my-counter", metrics.DayLabels())
```