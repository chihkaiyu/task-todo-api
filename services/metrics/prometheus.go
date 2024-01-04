package metric

import (
	"fmt"
	"sync"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type promMetric struct {
	namespace  string
	gauges     sync.Map
	histograms sync.Map
	counters   sync.Map
	mutex      sync.Mutex
}

type promTimeTracker struct {
	timer *prometheus.Timer
}

func (pm *promMetric) Gauge(metricName string, value float64, tags []Tag) {
	hashKey := hash(pm.namespace, metricName)
	labels := tagsToLabels(tags)

	pm.mutex.Lock()
	defer pm.mutex.Unlock()
	if collector, ok := pm.gauges.Load(hashKey); ok {
		gaugeVec := collector.(*prometheus.GaugeVec)
		gauge, err := gaugeVec.GetMetricWith(labels)
		if err != nil {
			log.Error().Err(err).
				Str("hashKey", hashKey).
				Str("namespace", pm.namespace).
				Msg("gaugeVec.GetMetricWith failed")
			return
		}
		gauge.Set(value)
		return
	}

	opts := prometheus.GaugeOpts{
		Namespace: pm.namespace,
		Name:      metricName,
	}

	keyArr, _ := tagsToKeyValueArray(tags)
	gaugeVec := prometheus.NewGaugeVec(opts, keyArr)
	if err := prometheus.Register(gaugeVec); err != nil {
		log.Error().Err(err).
			Str("hashKey", hashKey).
			Array("labels", keyArr).
			Str("namespace", pm.namespace).
			Msg("prometheus.Register failed")
		return
	}

	pm.gauges.Store(hashKey, gaugeVec)
	gauge, err := gaugeVec.GetMetricWith(labels)
	if err != nil {
		log.Error().Err(err).
			Str("hashKey", hashKey).
			Str("namespace", pm.namespace).
			Msg("gaugeVec.GetMetricWith failed")
		return
	}
	gauge.Set(value)
}

func (pm *promMetric) Time(metricName string, tags []Tag) Ender {
	hashKey := hash(pm.namespace, metricName)
	labels := tagsToLabels(tags)

	pm.mutex.Lock()
	defer pm.mutex.Unlock()
	if collector, ok := pm.histograms.Load(hashKey); ok {
		histogramVec := collector.(*prometheus.HistogramVec)
		histogram, err := histogramVec.GetMetricWith(labels)
		if err != nil {
			log.Error().Err(err).
				Str("hashKey", hashKey).
				Str("namespace", pm.namespace).
				Msg("prometheus.Register failed")
			return &fakeEnd{}
		}
		timer := prometheus.NewTimer(histogram)

		return &promTimeTracker{
			timer: timer,
		}
	}

	opts := prometheus.HistogramOpts{
		Namespace: pm.namespace,
		Name:      metricName,
	}

	keyArr, _ := tagsToKeyValueArray(tags)
	histogramVec := prometheus.NewHistogramVec(opts, keyArr)
	if err := prometheus.Register(histogramVec); err != nil {
		log.Error().Err(err).
			Str("hashKey", hashKey).
			Str("namespace", pm.namespace).
			Msg("prometheus.Register failed")
		return &fakeEnd{}
	}

	pm.histograms.Store(hashKey, histogramVec)
	histogram, err := histogramVec.GetMetricWith(labels)
	if err != nil {
		log.Error().Err(err).
			Str("hashKey", hashKey).
			Str("namespace", pm.namespace).
			Msg("prometheus.Register failed")
		return &fakeEnd{}
	}
	timer := prometheus.NewTimer(histogram)

	return &promTimeTracker{
		timer: timer,
	}
}

func (pt *promTimeTracker) End() {
	pt.timer.ObserveDuration()
}

func (pm *promMetric) Counter(metricName string, value float64, tags []Tag) {
	hashKey := hash(pm.namespace, metricName)
	labels := tagsToLabels(tags)

	pm.mutex.Lock()
	defer pm.mutex.Unlock()
	if collector, ok := pm.counters.Load(hashKey); ok {
		counterVec := collector.(*prometheus.CounterVec)
		counter, err := counterVec.GetMetricWith(labels)
		if err != nil {
			log.Error().Err(err).
				Str("hashKey", hashKey).
				Str("namespace", pm.namespace).
				Msg("counterVec.GetMetricWith failed")
			return
		}
		counter.Add(value)
		return
	}

	opts := prometheus.CounterOpts{
		Namespace: pm.namespace,
		Name:      metricName,
	}

	keyArr, _ := tagsToKeyValueArray(tags)
	counterVec := prometheus.NewCounterVec(opts, keyArr)
	if err := prometheus.Register(counterVec); err != nil {
		log.Error().Err(err).
			Str("hashKey", hashKey).
			Str("namespace", pm.namespace).
			Msg("prometheus.Register failed")
		return
	}

	pm.counters.Store(hashKey, counterVec)
	counter, err := counterVec.GetMetricWith(labels)
	if err != nil {
		log.Error().Err(err).
			Str("hashKey", hashKey).
			Str("namespace", pm.namespace).
			Msg("counterVec.GetMetricWith failed")
		return
	}
	counter.Add(value)
}

func hash(namespace, metricName string) string {
	return fmt.Sprintf("%s:%s", namespace, metricName)
}

type tagArray []string

func (ta tagArray) MarshalZerologArray(a *zerolog.Array) {
	for _, t := range ta {
		a.Str(t)
	}
}

func tagsToKeyValueArray(tags []Tag) (tagArray, tagArray) {
	key := make([]string, len(tags))
	value := make([]string, len(tags))
	for i := 0; i < len(tags); i++ {
		key[i] = tags[i].Name
		value[i] = tags[i].Value
	}

	return key, value
}

func tagsToLabels(tags []Tag) prometheus.Labels {
	labels := prometheus.Labels{}
	for _, t := range tags {
		labels[t.Name] = t.Value
	}

	return labels
}
