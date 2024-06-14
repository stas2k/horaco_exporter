package clients

import (
	"encoding/hex"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestHoracoClient(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(mockMainHandler))
	defer ts.Close()

	tc := ts.Client()
	hash_a, _ := hex.DecodeString("f6fdffe48c908deb0f4c3bd36c032e72")

	scraper := NewHoracoClient(ts.URL, ([16]byte)(hash_a))
	scraper.h_client = tc

	t.Run("GetPortStats", func(t *testing.T) {
		out, err := scraper.GetPortStats()
		if err != nil {
			t.Error(err)
		}
		if !reflect.DeepEqual(out,
			[]portStats{
				{
					State:      true,
					LinkStatus: true,
					PktCount: packetStats{
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
					PktCount: packetStats{
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
					PktCount: packetStats{
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
					PktCount: packetStats{
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
					PktCount: packetStats{
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
					PktCount: packetStats{
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
					PktCount: packetStats{
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
					PktCount: packetStats{
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
					PktCount: packetStats{
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
			}) {
			t.Errorf("does not equal expected output")
		}
	})

	t.Run("GetSystemInfo", func(t *testing.T) {
		out, err := scraper.GetSystemInfo()
		if err != nil {
			t.Error(err)
		}
		if !reflect.DeepEqual(out,
			&systemInfo{
				Model:           "WAMJHJ-8125MNG",
				MacAddress:      "1C:2A:12:34:56:78",
				FirmwareVersion: "V1.9",
				FirmwareDate:    "Jan 03 2024",
				HardwareVersion: "V1.1",
			}) {
			t.Errorf("does not equal expected output")
		}
	})
}
