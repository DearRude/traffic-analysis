package main

import (
	"github.com/prometheus/client_golang/prometheus"

	"sync"
)

var (
	// congestion to value map
	convertCongestionMap = map[string]int{
		"low":      0,
		"moderate": 1,
		"heavy":    2,
		"severe":   3,
	}

	congesMap      map[string]map[string]int
	congesCountMap map[string]map[string]int

	mtx sync.Mutex
)

type metrics struct {
	congestionPercent   *prometheus.GaugeVec
	processedTilesTotal *prometheus.CounterVec
}

func NewMetrics(reg prometheus.Registerer) *metrics {
	m := &metrics{
		congestionPercent: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name: "traffic_congestion_percent",
			Help: "Percentage of traffic congestion of each city for each road_class",
		}, []string{"city", "road_class"}),
		processedTilesTotal: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "traffic_processed_total",
				Help: "Number processed traffic of each city for each road_class",
			}, []string{"city", "road_class"}),
	}
	reg.MustRegister(m.congestionPercent)
	reg.MustRegister(m.processedTilesTotal)
	return m
}

func writeCongestionMetrics(ways []Way, traffics []Traffic) {
	mtx.Lock()
	defer mtx.Unlock()

	for i := range traffics {
		if congesMap[ways[i].City] == nil {
			congesMap[ways[i].City] = make(map[string]int)
		}
		congesMap[ways[i].City][ways[i].RoadClass] += convertCongestionMap[traffics[i].TrafficState]

		if congesCountMap[ways[i].City] == nil {
			congesCountMap[ways[i].City] = make(map[string]int)
		}
		congesCountMap[ways[i].City][ways[i].RoadClass] += 1
	}
}

func setCongestionMetrics(m *metrics) {
	for city, values := range congesMap {
		for roadClass, congestion := range values {
			count := congesCountMap[city][roadClass]
			percentRatio := float64((100 / 3) / count)
			m.congestionPercent.With(prometheus.Labels{"city": city, "road_class": roadClass}).Set(float64(congestion) * percentRatio)
			m.processedTilesTotal.With(prometheus.Labels{"city": city, "road_class": roadClass}).Add(float64(count))
		}
	}
}

func resetCongestionMetrics() {
	congesMap = make(map[string]map[string]int)
	congesCountMap = make(map[string]map[string]int)
}
