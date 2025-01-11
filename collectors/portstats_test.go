package collectors

import (
	"testing"
	"github.com/prometheus/client_golang/prometheus"

	"github.com/stas2k/horaco_exporter/clients"
	"github.com/stas2k/horaco_exporter/testutil"
)

type mockStatsClient bool

func (c *mockStatsClient) GetPortStats() ([]clients.PortStats, error) {
	return []clients.PortStats{
				{
					State:      true,
					LinkStatus: true,
					PktCount: clients.PacketStats{
						TxGood: 1,
						TxBad:  11,
						RxGood: 111,
						RxBad:  1111,
					},
					LinkSpeedSet:         0,
					LinkSpeedActual:      1000,
					LinkFullDuplexSet:    true,
					LinkFullDuplexActual: true,
					FlowControlSet:       false,
					FlowControlActual:    false,
				},
				{
					State:      true,
					LinkStatus: true,
					PktCount: clients.PacketStats{
						TxGood: 2,
						TxBad:  22,
						RxGood: 222,
						RxBad:  2222,
					},
					LinkSpeedSet:         100,
					LinkSpeedActual:      100,
					LinkFullDuplexSet:    true,
					LinkFullDuplexActual: true,
					FlowControlSet:       true,
					FlowControlActual:    true,
				},

				{
					State:      false,
					LinkStatus: true,
					PktCount: clients.PacketStats{
						TxGood: 3,
						TxBad:  33,
						RxGood: 333,
						RxBad:  3333,
					},
					LinkSpeedSet:         100,
					LinkSpeedActual:      100,
					LinkFullDuplexSet:    false,
					LinkFullDuplexActual: false,
					FlowControlSet:       true,
					FlowControlActual:    false,
				},
				{
					State:      true,
					LinkStatus: true,
					PktCount: clients.PacketStats{
						TxGood: 4,
						TxBad:  44,
						RxGood: 444,
						RxBad:  4444,
					},
					LinkSpeedSet:         10,
					LinkSpeedActual:      10,
					LinkFullDuplexSet:    true,
					LinkFullDuplexActual: true,
					FlowControlSet:       false,
					FlowControlActual:    false,
				},
				{
					State:      true,
					LinkStatus: true,
					PktCount: clients.PacketStats{
						TxGood: 5,
						TxBad:  55,
						RxGood: 555,
						RxBad:  5555,
					},
					LinkSpeedSet:         10,
					LinkSpeedActual:      10,
					LinkFullDuplexSet:    false,
					LinkFullDuplexActual: false,
					FlowControlSet:       false,
					FlowControlActual:    false,
				},
				{
					State:      true,
					LinkStatus: true,
					PktCount: clients.PacketStats{
						TxGood: 6,
						TxBad:  66,
						RxGood: 666,
						RxBad:  6666,
					},
					LinkSpeedSet:         2500,
					LinkSpeedActual:      2500,
					LinkFullDuplexSet:    true,
					LinkFullDuplexActual: true,
					FlowControlSet:       true,
					FlowControlActual:    true,
				},
				{
					State:      true,
					LinkStatus: false,
					PktCount: clients.PacketStats{
						TxGood: 7,
						TxBad:  77,
						RxGood: 777,
						RxBad:  7777,
					},
					LinkSpeedSet:         0,
					LinkSpeedActual:      0,
					LinkFullDuplexSet:    true,
					LinkFullDuplexActual: false,
					FlowControlSet:       true,
					FlowControlActual:    false,
				},
				{
					State:      false,
					LinkStatus: false,
					PktCount: clients.PacketStats{
						TxGood: 8,
						TxBad:  88,
						RxGood: 888,
						RxBad:  8888,
					},
					LinkSpeedSet:         0,
					LinkSpeedActual:      0,
					LinkFullDuplexSet:    true,
					LinkFullDuplexActual: false,
					FlowControlSet:       false,
					FlowControlActual:    false,
				},
				{
					State:      true,
					LinkStatus: true,
					PktCount: clients.PacketStats{
						TxGood: 9,
						TxBad:  99,
						RxGood: 999,
						RxBad:  9999,
					},
					LinkSpeedSet:         10000,
					LinkSpeedActual:      10000,
					LinkFullDuplexSet:    true,
					LinkFullDuplexActual: true,
					FlowControlSet:       true,
					FlowControlActual:    false,
				},
	}, nil
}

func TestPortStatsCollector(t *testing.T) {
	client := new(mockStatsClient)
	stats_collector := NewPortStatsCollector("horaco_exporter", client)

	registry := prometheus.NewRegistry()
	registry.MustRegister(stats_collector)

	t.Run("GetPortStats", func(t *testing.T) {
		out, err := registry.Gather()
		if err != nil {
			t.Fatalf("Error gathering metrics: %v", err)
		}

		ok := testutil.CompareMetrics(out, testutil.MetricTests{
			"horaco_exporter_probe_success": {
				Labels: testutil.LabelTests{
					"probe": "port",
				},
				Value: 1.0,
			},
		})
		if !ok {
			t.Fatalf("Unexpected response from collector: %v", out)
		}
		ok = testutil.CompareMetrics(out, testutil.MetricTests{
			"horaco_exporter_probe_success": {
				Labels: testutil.LabelTests{
					"probe": "info",
				},
				Value: 1.0,
			},
		})
		if !ok {
			t.Fatalf("Unexpected response from collector: %v", out)
		}
	})
}
