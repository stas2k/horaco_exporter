package clients

import (
	"time"
	"bufio"
	"encoding/hex"
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"strings"
)

type packetStats struct {
	TxGood, TxBad, RxGood, RxBad int
}
type portStats struct {
	// from monitoring page
	State      bool
	LinkStatus bool
	PktCount   packetStats
	// from config page
	LinkSpeedSet    int
	LinkSpeedActual int

	LinkFullDuplexSet    bool
	LinkFullDuplexActual bool

	FlowControlSet    bool
	FlowControlActual bool
}

// node_openwrt_info{board_name="xiaomi,mi-router-3g", id="OpenWrt", model="Xiaomi Mi Router 3G", release="23.05.3", revision="r23809-234f1a2efa", system="MediaTek MT7621 ver:1 eco:3", target="ramips/mt7621"}
type systemInfo struct {
	Model           string
	MacAddress      string
	FirmwareVersion string
	FirmwareDate    string
	HardwareVersion string
}

type HoracoClient struct {
	header   *http.Header
	line_reg *regexp.Regexp
	h_client *http.Client
	base_url string
}

const PORT_URL = "/port.cgi"
const STATS_URL = PORT_URL + "?page=stats"
const INFO_URL = "/info.cgi"

func (client *HoracoClient) parseField(s *bufio.Scanner) (string, error) {
	s.Scan()
	match := client.line_reg.FindStringSubmatch(s.Text())
	if match == nil {
		return "", errors.New("empty field")
	}
	return match[1], nil
}

func (client *HoracoClient) parseFieldNum(s *bufio.Scanner) (int, error) {
	line, err := client.parseField(s)
	if err != nil {
		return 0, err
	}
	cnt, err := strconv.Atoi(line)
	if err != nil {
		return 0, err
	}
	return cnt, nil
}

func (client *HoracoClient) getURL(p string) (*http.Response, error) {
	url := fmt.Sprintf(client.base_url + p)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header = *client.header
	res, err := client.h_client.Do(req)
	if err != nil {
		return nil, err
	}
	if (res.StatusCode < 100) || (res.StatusCode > 299) {
		return nil, fmt.Errorf("Unexpected HTTP response: %s", res.Status)
	}
	return res, nil
}

func (client *HoracoClient) GetSystemInfo() (*systemInfo, error) {
	info_resp, err := client.getURL(INFO_URL)
	if err != nil {
		return nil, fmt.Errorf("error getting INFO_URL: %w", err)
	}
	defer info_resp.Body.Close()

	info_scan := bufio.NewScanner(info_resp.Body)
	for info_scan.Scan() {
		if strings.HasSuffix(info_scan.Text(), "<table>") {
			break
		}
	}
	if err := info_scan.Err(); err != nil {
		return nil, err
	}

	var info systemInfo

	info_scan.Scan()
	line := info_scan.Text()
	if !strings.HasSuffix(line, "<tr>") {
		return nil, errors.New("unexpected line when waiting for start of table row: " + line)
	}

	info_scan.Scan() // skip line with header
	line, err = client.parseField(info_scan)
	if err != nil {
		return nil, fmt.Errorf("error scanning Model: %w", err)
	}
	info.Model = line
	info_scan.Scan() // skip line with header
	info_scan.Scan() // skip line with header

	info_scan.Scan() // skip line with header
	line, err = client.parseField(info_scan)
	if err != nil {
		return nil, fmt.Errorf("error scanning MacAddress: %w", err)
	}
	info.MacAddress = line
	info_scan.Scan() // skip line with header
	info_scan.Scan() // skip line with header

	info_scan.Scan() // skip IP
	info_scan.Scan()
	info_scan.Scan()
	info_scan.Scan()

	info_scan.Scan() // skip Netmask
	info_scan.Scan()
	info_scan.Scan()
	info_scan.Scan()

	info_scan.Scan() // skip Gateway
	info_scan.Scan()
	info_scan.Scan()
	info_scan.Scan()

	info_scan.Scan() // skip line with header
	line, err = client.parseField(info_scan)
	if err != nil {
		return nil, fmt.Errorf("error scanning FirmwareVersion: %w", err)
	}
	info.FirmwareVersion = line
	info_scan.Scan() // skip line with header
	info_scan.Scan() // skip line with header

	info_scan.Scan() // skip line with header
	line, err = client.parseField(info_scan)
	if err != nil {
		return nil, fmt.Errorf("error scanning FirmwareDate: %w", err)
	}
	info.FirmwareDate = line
	info_scan.Scan() // skip line with header
	info_scan.Scan() // skip line with header

	info_scan.Scan() // skip line with header
	line, err = client.parseField(info_scan)
	if err != nil {
		return nil, fmt.Errorf("error scanning HardwareVersion: %w", err)
	}
	info.HardwareVersion = line
	info_scan.Scan() // skip line with header
	info_scan.Scan() // skip line with header

	return &info, nil
}

func (client *HoracoClient) GetPortStats() ([]portStats, error) {

	stat_resp, err := client.getURL(STATS_URL)
	if err != nil {
		return nil, fmt.Errorf("error getting STATS_URL: %w", err)
	}
	stat_scan := bufio.NewScanner(stat_resp.Body)
	for stat_scan.Scan() {
		if strings.HasSuffix(stat_scan.Text(), "</tr>") {
			break
		}
	}
	if err := stat_scan.Err(); err != nil {
		return nil, err
	}

	// Get Port Stats
	ports := make([]portStats, 9)

	for i := 0; i < len(ports); i++ {
		stat_scan.Scan()
		line := stat_scan.Text()
		if !strings.HasSuffix(line, "<tr>") {
			return nil, errors.New("unexpected line when waiting for start of table row: " + line)
		}

		ps := &ports[i]

		stat_scan.Scan() // skip line with Port Number

		line, err = client.parseField(stat_scan)
		if err != nil {
			return nil, fmt.Errorf("error scanning State: %w", err)
		}
		if line == "Enable" {
			ps.State = true
		} else if line == "Disable" {
			ps.State = false
		} else {
			return nil, fmt.Errorf("error scanning State: unexpected value: %s", line)
		}

		line, err = client.parseField(stat_scan)
		if err != nil {
			return nil, fmt.Errorf("error scanning LinkStatus: %w", err)
		}
		if line == "Link Up" {
			ps.LinkStatus = true
		} else if line == "Link Down" {
			ps.LinkStatus = false
		} else {
			return nil, fmt.Errorf("error scanning LinkStatus: unexpected value: %s", line)
		}

		cnt, err := client.parseFieldNum(stat_scan)
		if err != nil {
			return nil, fmt.Errorf("error scanning TxGood: %w", err)
		}
		ps.PktCount.TxGood = cnt
		cnt, err = client.parseFieldNum(stat_scan)
		if err != nil {
			return nil, fmt.Errorf("error scanning TxBad: %w", err)
		}
		ps.PktCount.TxBad = cnt
		cnt, err = client.parseFieldNum(stat_scan)
		if err != nil {
			return nil, fmt.Errorf("error scanning RxGood: %w", err)
		}
		ps.PktCount.RxGood = cnt
		cnt, err = client.parseFieldNum(stat_scan)
		if err != nil {
			return nil, fmt.Errorf("error scanning RxBad: %w", err)
		}
		ps.PktCount.RxBad = cnt

		stat_scan.Scan() // skip line with </tr>
	}

	stat_resp.Body.Close()

	//////////////////////////////////////////
	/*
	   Config:
	   2.5G Full
	   1000Full
	   100 Full
	   100 Half
	   10 Full
	   10 Half
	   10G Full
	   Actual:
	   1000
	   2500
	   5000
	   Half
	   Full
	*/
	port_resp, err := client.getURL(PORT_URL)
	if err != nil {
		return nil, fmt.Errorf("error getting PORT_URL: %w", err)
	}
	port_scan := bufio.NewScanner(port_resp.Body)
	for port_scan.Scan() {
		if port_scan.Text() == "<table border=\"1\">" {
			break
		}
	}
	for port_scan.Scan() {
		if strings.HasSuffix(port_scan.Text(), "Actual</th>") {
			break
		}
	}
	port_scan.Scan() // skip line with </tr>
	port_scan.Scan() // skip line with </tr>
	port_scan.Scan() // skip line with </tr>
	if err := port_scan.Err(); err != nil {
		return nil, err
	}

	for i := 0; i < len(ports); i++ {
		port_scan.Scan()
		line := port_scan.Text()
		if !strings.HasSuffix(line, "<tr>") {
			return nil, errors.New("unexpected line when waiting for start of table row: " + line)
		}
		// Skip 2 cells - Port number and State
		port_scan.Scan()
		port_scan.Scan()

		ps := &ports[i]

		line, err = client.parseField(port_scan)
		if err != nil {
			return nil, fmt.Errorf("error scanning LinkSpeedSet: %w", err)
		}
		switch line {
		case "Auto":
			ps.LinkSpeedSet = 0
			ps.LinkFullDuplexSet = true
		case "10G Full":
			ps.LinkSpeedSet = 10000
			ps.LinkFullDuplexSet = true
		case "2.5G Full":
			ps.LinkSpeedSet = 2500
			ps.LinkFullDuplexSet = true
		case "1000Full":
			ps.LinkSpeedSet = 1000
			ps.LinkFullDuplexSet = true
		case "100 Full":
			ps.LinkSpeedSet = 100
			ps.LinkFullDuplexSet = true
		case "100 Half":
			ps.LinkSpeedSet = 100
			ps.LinkFullDuplexSet = false
		case "10 Full":
			ps.LinkSpeedSet = 10
			ps.LinkFullDuplexSet = true
		case "10 Half":
			ps.LinkSpeedSet = 10
			ps.LinkFullDuplexSet = false
		default:
			return nil, fmt.Errorf("error scanning LinkSpeedSet: unexpected value: %s", line)
		}

		line, err = client.parseField(port_scan)
		if err != nil {
			return nil, fmt.Errorf("error scanning LinkSpeedActual: %w", err)
		}
		switch line {
		case "Link Down":
			ps.LinkSpeedActual = 0
			ps.LinkFullDuplexActual = false
		case "10G Full":
			ps.LinkSpeedActual = 10000
			ps.LinkFullDuplexActual = true
		case "2500Full":
			ps.LinkSpeedActual = 2500
			ps.LinkFullDuplexActual = true
		case "1000Full":
			ps.LinkSpeedActual = 1000
			ps.LinkFullDuplexActual = true
		case "100 Full":
			ps.LinkSpeedActual = 100
			ps.LinkFullDuplexActual = true
		case "100 Half":
			ps.LinkSpeedActual = 100
			ps.LinkFullDuplexActual = false
		case "10 Full":
			ps.LinkSpeedActual = 10
			ps.LinkFullDuplexActual = true
		case "10 Half":
			ps.LinkSpeedActual = 10
			ps.LinkFullDuplexActual = false
		default:
			return nil, fmt.Errorf("error scanning LinkSpeedActual: unexpected value: %s", line)
		}

		line, err = client.parseField(port_scan)
		if err != nil {
			return nil, fmt.Errorf("error scanning FlowControlSet: %w", err)
		}
		if line == "On" {
			ps.FlowControlSet = true
		} else if line == "Off" {
			ps.FlowControlSet = false
		} else {
			return nil, fmt.Errorf("error scanning FlowControlSet: unexpected value: %s", line)
		}

		line, err = client.parseField(port_scan)
		if err != nil {
			return nil, fmt.Errorf("error scanning FlowControlActual: %w", err)
		}
		if line == "On" {
			ps.FlowControlActual = true
		} else if line == "Off" {
			ps.FlowControlActual = false
		} else {
			return nil, fmt.Errorf("error scanning FlowControlActual: unexpected value: %s", line)
		}
		port_scan.Scan() // skip line with </tr>

	}

	port_resp.Body.Close()

	return ports, nil
}

func NewHoracoClient(base_url string, hash [16]byte) *HoracoClient {
	return &HoracoClient{
		base_url: base_url,
		line_reg: regexp.MustCompile(`\s+<td(?: [^>]+)?>([^<+]+)</td>`),
		header: &http.Header{
			"Cookie": {"admin=" + hex.EncodeToString(hash[:])},
		},
		h_client: &http.Client{Timeout: 3 * time.Second},
	}
}
