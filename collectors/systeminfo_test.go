package collectors

import (
	"testing"
	"github.com/prometheus/client_golang/prometheus"

	"github.com/stas2k/horaco_exporter/clients"
	"github.com/stas2k/horaco_exporter/testutil"
)

type mockInfoClient bool

func (c *mockInfoClient) GetSystemInfo() (*clients.SystemInfo, error) {
	return &clients.SystemInfo{
		Model:           "WAMJHJ-8125MNG",
		MacAddress:      "1C:2A:12:34:56:78",
		FirmwareVersion: "V1.9",
		FirmwareDate:    "Jan 03 2024",
		HardwareVersion: "V1.1",
	}, nil
}

func TestSystemInfoCollector(t *testing.T) {
	client := new(mockInfoClient)
	info_collector := NewSystemInfoCollector("horaco_exporter", client)

	registry := prometheus.NewRegistry()
	registry.MustRegister(info_collector)

	t.Run("GetSystemInfo", func(t *testing.T) {
		out, err := registry.Gather()
		if err != nil {
			t.Fatalf("Error gathering metrics: %v", err)
		}
		ok := testutil.CompareMetrics(out, testutil.MetricTests{
			"horaco_exporter_probe_success": {
				Labels: testutil.LabelTests{
					"probe": "info",
				},
				Value: 1.0,
			},
			"horaco_exporter_switch_info": {
				Labels: testutil.LabelTests{
					"firmware_date": "Jan 03 2024",
					"firmware_ver":  "V1.9",
					"hardware_ver":  "V1.1",
					"mac":           "1C:2A:12:34:56:78",
					"model":         "WAMJHJ-8125MNG",
				},
				Value: 1.0,
			},
		})
		if !ok {
			t.Fatalf("Unexpected response from collector: %v", out)
		}
	})
}
