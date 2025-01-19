package collectors

import (
	"log"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/stas2k/horaco_exporter/clients"
)

type InfoClient interface {
	GetSystemInfo() (*clients.SystemInfo, error)
}

type SystemInfoCollector struct {
	namespace  string
	client     InfoClient
	infoMetric *prometheus.GaugeVec
}

func NewSystemInfoCollector(namespace string, client InfoClient) *SystemInfoCollector {
	return &SystemInfoCollector{
		namespace: namespace,
		client:    client,
		infoMetric: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Subsystem: "switch",
				Name:      "info",
				Help:      "A metric with a constant '1' value labeled by system information scraped from the switch",
			},
			[]string{
				"model",
				"mac",
				"firmware_ver",
				"firmware_date",
				"hardware_ver",
			},
		),
	}
}

func (c *SystemInfoCollector) Collect(ch chan<- prometheus.Metric) {
	success := prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: c.namespace,
			Name:      "probe_success",
			Help:      "Metric that exposes whether underlying requests to device succeeded. Value '1' request was succesful.",
		},
		[]string{"probe"},
	).With(prometheus.Labels{"probe": "info"})

	info, err := c.client.GetSystemInfo()
	if err != nil {
		log.Printf("Error collecting system info metrics: %s", err)
		success.Set(0.0)
		success.Collect(ch)
		return
	}
	success.Set(1.0)
	success.Collect(ch)

	labels := prometheus.Labels{
		"model":         info.Model,
		"mac":           info.MacAddress,
		"firmware_ver":  info.FirmwareVersion,
		"firmware_date": info.FirmwareDate,
		"hardware_ver":  info.HardwareVersion,
	}

	c.infoMetric.With(labels).Set(1.0)
	c.infoMetric.Collect(ch)
}
func (c *SystemInfoCollector) Describe(ch chan<- *prometheus.Desc) {
	c.infoMetric.Describe(ch)
}
