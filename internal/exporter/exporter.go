package exporter

import (
	"github.com/alexandrevilain/atome_exporter/pkg/atome"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/log"
	"github.com/sirupsen/logrus"
)

const (
	namespace = "atome"
)

// Exporter collects Atome statistics from the given atome client
type Exporter struct {
	logger      *logrus.Logger
	atome       *atome.Client
	up          *prometheus.Desc
	consumption *prometheus.Desc
	price       *prometheus.Desc
	co2impact   *prometheus.Desc
}

// New creates a new instance of the atome exporter
func New(logger *logrus.Logger, atome *atome.Client) *Exporter {
	return &Exporter{
		logger: logger,
		atome:  atome,
		up: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "up"),
			"Was the last query of atome successful.",
			nil, nil,
		),
		consumption: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "consumption"),
			"The current energy consumption in wh",
			nil, nil,
		),
		price: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "price"),
			"The current price it costs",
			nil, nil,
		),
		co2impact: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "co2impact"),
			"The current impact in co2 for the consumption",
			nil, nil,
		),
	}
}

// Describe describes all the metrics exported by this exporter
func (e *Exporter) Describe(ch chan<- *prometheus.Desc) {
	ch <- e.up
	ch <- e.consumption
	ch <- e.price
	ch <- e.co2impact
}

// Collect fetches the stats from the atome client
func (e *Exporter) Collect(ch chan<- prometheus.Metric) {
	e.logger.Infoln("Collecting metrics")

	up := 1
	consumption, err := e.atome.RetriveDayConsumption()
	if err != nil {
		log.Error(err)
		up = 0
	}

	ch <- prometheus.MustNewConstMetric(e.up, prometheus.GaugeValue, float64(up))
	ch <- prometheus.MustNewConstMetric(e.consumption, prometheus.GaugeValue, float64(consumption.Total))
	ch <- prometheus.MustNewConstMetric(e.price, prometheus.GaugeValue, consumption.Price)
	ch <- prometheus.MustNewConstMetric(e.co2impact, prometheus.GaugeValue, float64(consumption.Co2Impact))
}
