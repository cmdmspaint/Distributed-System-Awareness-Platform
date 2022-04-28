package metrics

import (
	"Distributed-System-Awareness-Platform/src/models"
	"github.com/prometheus/client_golang/prometheus"
)

func CreateMetrics(ss []*models.LogStrategy) map[string]*prometheus.GaugeVec {

	mmap := make(map[string]*prometheus.GaugeVec)
	for _, s := range ss {
		labels := []string{}

		for k := range s.Tags {
			labels = append(labels, k)
		}
		m := prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name: s.MetricName,
			Help: s.MetricHelp,
		}, labels)
		mmap[s.MetricName] = m

	}
	return mmap
}
