package collectors

import (
	"fmt"
	"log"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/stas2k/horaco_exporter/clients"
)

type PortClient interface {
	GetPortStats() ([]clients.PortStats, error)
}

type PortStatsCollector struct {
	namespace string
	client    PortClient
	gauges    []prometheus.GaugeVec
	counters  []prometheus.CounterVec
}

func NewPortStatsCollector(namespace string, client PortClient) *PortStatsCollector {
	return &PortStatsCollector{
		client:    client,
		namespace: namespace,
	}
}

func (c *PortStatsCollector) Collect(ch chan<- prometheus.Metric) {
	success := prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: c.namespace,
			Name:      "probe_success",
			Help:      "Metric that exposes whether underlying requests to device succeeded. Value '1' request was succesful.",
		},
		[]string{"probe"},
	).With(prometheus.Labels{"probe": "port"})
	stats, err := c.client.GetPortStats()
	if err != nil {
		log.Printf("Error collecting port metrics: %s", err)
		success.Set(0.0)
		success.Collect(ch)
		return
	}
	success.Set(1.0)
	success.Collect(ch)

	for i := 0; i < len(stats); i++ {
		port_labels := []string{
			"device",
		}
		port := fmt.Sprintf("port%v", i+1)
		stat := stats[i]

		m := prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: c.namespace,
				Subsystem: "port",
				Name:      "state",
				Help:      "Metric that exposes administrator set link state. Value '1' means the link is enabled.",
			},
			port_labels,
		)
		if stat.State {
			m.With(prometheus.Labels{"device": port}).Set(1.0)
		} else {
			m.With(prometheus.Labels{"device": port}).Set(0.0)
		}
		c.gauges = append(c.gauges, *m)

		m = prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: c.namespace,
				Subsystem: "port",
				Name:      "carrier",
				Help:      "Metric that exposes actual link carrier state. Value '1' means the link is up.",
			},
			port_labels,
		)
		if stat.LinkStatus {
			m.With(prometheus.Labels{"device": port}).Set(1.0)
		} else {
			m.With(prometheus.Labels{"device": port}).Set(0.0)
		}
		c.gauges = append(c.gauges, *m)

		m = prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: c.namespace,
				Subsystem: "port",
				Name:      "speed_bytes_set",
				Help:      "Metric that exposes administrator set desired link speed in bytes per second. Value 0 means auto negotiation",
			},
			port_labels,
		)
		m.With(prometheus.Labels{"device": port}).Set(float64(stat.LinkSpeedSet * 125000))
		c.gauges = append(c.gauges, *m)

		m = prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: c.namespace,
				Subsystem: "port",
				Name:      "speed_bytes",
				Help:      "Metric that exposes negotiated link speed in bytes per second",
			},
			port_labels,
		)
		m.With(prometheus.Labels{"device": port}).Set(float64(stat.LinkSpeedActual * 125000))
		c.gauges = append(c.gauges, *m)

		m = prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: c.namespace,
				Subsystem: "port",
				Name:      "duplex_set",
				Help:      "Metric that exposes administrator set duplex state. Value '1' means full duplex is desired.",
			},
			port_labels,
		)
		if stat.LinkFullDuplexSet {
			m.With(prometheus.Labels{"device": port}).Set(1.0)
		} else {
			m.With(prometheus.Labels{"device": port}).Set(0.0)
		}
		c.gauges = append(c.gauges, *m)

		m = prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: c.namespace,
				Subsystem: "port",
				Name:      "duplex",
				Help:      "Metric that exposes negotiated duplex state. Value '1' means full duplex negotiated.",
			},
			port_labels,
		)
		if stat.LinkFullDuplexActual {
			m.With(prometheus.Labels{"device": port}).Set(1.0)
		} else {
			m.With(prometheus.Labels{"device": port}).Set(0.0)
		}
		c.gauges = append(c.gauges, *m)

		m = prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: c.namespace,
				Subsystem: "port",
				Name:      "flow_control_set",
				Help:      "Metric that exposes administrator set flow control state. Value '1' means flow control is desired.",
			},
			port_labels,
		)
		if stat.FlowControlSet {
			m.With(prometheus.Labels{"device": port}).Set(1.0)
		} else {
			m.With(prometheus.Labels{"device": port}).Set(0.0)
		}
		c.gauges = append(c.gauges, *m)

		m = prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: c.namespace,
				Subsystem: "port",
				Name:      "flow_control",
				Help:      "Metric that exposes negotiated flow control state. Value '1' means flow control is negotiated.",
			},
			port_labels,
		)
		if stat.FlowControlActual {
			m.With(prometheus.Labels{"device": port}).Set(1.0)
		} else {
			m.With(prometheus.Labels{"device": port}).Set(0.0)
		}
		c.gauges = append(c.gauges, *m)

		cnt := prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: c.namespace,
				Subsystem: "port",
				Name:      "transmit_frames_total",
				Help:      "Metric that exposes number of successfuly transmitted frames.",
			},
			port_labels,
		)
		cnt.With(prometheus.Labels{"device": port}).Add(float64(stat.PktCount.TxGood))
		c.counters = append(c.counters, *cnt)

		cnt = prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: c.namespace,
				Subsystem: "port",
				Name:      "transmit_errs_total",
				Help:      "Metric that exposes number of errors transmitting frames.",
			},
			port_labels,
		)
		cnt.With(prometheus.Labels{"device": port}).Add(float64(stat.PktCount.TxBad))
		c.counters = append(c.counters, *cnt)

		cnt = prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: c.namespace,
				Subsystem: "port",
				Name:      "receive_frames_total",
				Help:      "Metric that exposes number of successfuly received frames.",
			},
			port_labels,
		)
		cnt.With(prometheus.Labels{"device": port}).Add(float64(stat.PktCount.RxGood))
		c.counters = append(c.counters, *cnt)

		cnt = prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: c.namespace,
				Subsystem: "port",
				Name:      "receive_errs_total",
				Help:      "Metric that exposes number of errors receiving frames.",
			},
			port_labels,
		)
		cnt.With(prometheus.Labels{"device": port}).Add(float64(stat.PktCount.RxBad))
		c.counters = append(c.counters, *cnt)

	}

	for i := 0; i < len(c.gauges); i++ {
		c.gauges[i].Collect(ch)
	}
	for i := 0; i < len(c.counters); i++ {
		c.counters[i].Collect(ch)
	}
}
func (c *PortStatsCollector) Describe(ch chan<- *prometheus.Desc) {

	for i := 0; i < len(c.gauges); i++ {
		c.gauges[i].Describe(ch)
	}
	for i := 0; i < len(c.counters); i++ {
		c.counters[i].Describe(ch)
	}
}
