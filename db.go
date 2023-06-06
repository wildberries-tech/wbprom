package wbprom

import (
	"database/sql"
	"github.com/prometheus/client_golang/prometheus"
)

type DbMetrics interface {
	ReadStatsFromDB(s *sql.DB)
}

type dbMetrics struct {
	NbMaxConns      prometheus.Gauge
	NbOpenConns     prometheus.Gauge
	NbUsedConns     prometheus.Gauge
	WaitCount       prometheus.Gauge
	WaitDurationSec prometheus.Summary
}

func NewDbMetrics(ns, subsystem, dbName string) *dbMetrics {
	return &dbMetrics{
		NbMaxConns:      newGauge(ns, subsystem, "nb_max_conns", dbName),
		NbOpenConns:     newGauge(ns, subsystem, "nb_open_conns", dbName),
		NbUsedConns:     newGauge(ns, subsystem, "nb_used_conns", dbName),
		WaitCount:       newGauge(ns, subsystem, "wait_count", dbName),
		WaitDurationSec: newSummary(ns, subsystem, "wait_duration_sec", dbName),
	}
}

func newGauge(ns, subsystem, name, labelDb string) prometheus.Gauge {
	var labels prometheus.Labels
	if labelDb != "" {
		labels = map[string]string{
			"db": labelDb,
		}
	}

	g := prometheus.NewGauge(
		prometheus.GaugeOpts{
			Namespace:   ns,
			Subsystem:   subsystem,
			Name:        name,
			ConstLabels: labels,
		})
	prometheus.MustRegister(g)
	return g
}

func newSummary(ns, subsystem, name, labelDb string) prometheus.Summary {
	var labels prometheus.Labels
	if labelDb != "" {
		labels = map[string]string{
			"db": labelDb,
		}
	}

	s := prometheus.NewSummary(
		prometheus.SummaryOpts{
			Namespace:   ns,
			Subsystem:   subsystem,
			Name:        name,
			ConstLabels: labels,
		})
	prometheus.MustRegister(s)
	return s
}

// Метод безопасно вызывать на закрытой БД
func (d *dbMetrics) ReadStatsFromDB(s *sql.DB) {
	stats := s.Stats()

	d.NbMaxConns.Set(float64(stats.MaxOpenConnections))
	d.NbOpenConns.Set(float64(stats.OpenConnections))
	d.NbUsedConns.Set(float64(stats.InUse))
	d.WaitCount.Set(float64(stats.WaitCount))
	d.WaitDurationSec.Observe(stats.WaitDuration.Seconds())

}
