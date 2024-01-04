package metric

import "sync"

// one metric instance per namespace
var (
	metInstance = map[string]*met{}
	lock        = &sync.Mutex{}
)

type Tag struct {
	Name  string
	Value string
}

type Service interface {
	Gauge(metricName string, value float64, tags []Tag)
	Time(metricName string, tags []Tag) Ender
	Counter(metricName string, value float64, tags []Tag)
	// TODO:
	// Histogram
}

type Ender interface {
	End()
}

type timeTracker struct {
	promEnder Ender
}

type fakeEnd struct{}

func (e *fakeEnd) End() {
}

type met struct {
	prom *promMetric
}

func New(namespace string) Service {
	if metInstance[namespace] == nil {
		lock.Lock()
		defer lock.Unlock()
		if metInstance[namespace] == nil {
			metInstance[namespace] = &met{
				prom: &promMetric{
					namespace:  namespace,
					gauges:     sync.Map{},
					counters:   sync.Map{},
					histograms: sync.Map{},
					mutex:      sync.Mutex{},
				},
			}
		}
	}

	return metInstance[namespace]
}

func (m *met) Gauge(metricName string, value float64, tags []Tag) {
	m.prom.Gauge(metricName, value, tags)
}

func (m *met) Time(metricName string, tags []Tag) Ender {
	promEnder := m.prom.Time(metricName, tags)
	return &timeTracker{
		promEnder: promEnder,
	}
}

func (m *met) Counter(metricName string, value float64, tags []Tag) {
	m.prom.Counter(metricName, value, tags)
}

func (t *timeTracker) End() {
	t.promEnder.End()
}
