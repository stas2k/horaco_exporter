# Horaco Exporter

Horaco_exporter is a Prometheus metrics exporter capable of scraping statistics from 8 port 2.5G Copper + 1 10G SFP port switches built on RTL8373+RTL8224 chipset.

It works by scraping the HTML of the administration interface, as I know no other way to get statistics off the switch.

The following models are known to work:

  - WAMJHJ-8125MNG, firmware V1.9

## Install

Horaco_exporter is a Go application, which can be built by the user, or downloaded from Releases page on girhub.

## Run

You can run the exporter without any arguments:

```shell
$ horaco_exporter
```

Afterwards you can visit the exporter's web UI to test it on port 8088(default).
Curl can be used to test as follows:

```shell
$ curl -G \
  -d user=admin \
  -d password=admin \
  -d target=http://192.168.2.1\
  http://localhost:8088/metrics
```

However, the default mode is insecure, as Prometheus needs to present switche's admin user credentials. You can store credentials on the exporter's side, and use them as endpoint allowlist.

Fill in a `auth` file as follows for switches' defualt credentials:

```
http://192.168.2.1 admin admin
```

Afterwards, run exporter as follows:

```shell
$ horaco_exporter -hosts /path/to/auth
```

In this mode, only `target` needs to be specified when scraping.
Curl can be used to test as follows:

```shell
$ curl -G \
  -d target=http://192.168.2.1\
  http://localhost:8088/metrics
```

## Scraping with Prometheus

The exporter is scraped in a simmilar manner to a blackbox_exporter.

The following configuration is a good starting point for allowlist mode:

```yaml
- job_name: horaco_exporter
  relabel_configs:
  - source_labels:
    - __address__
    target_label: __param_target
  - source_labels:
    - __param_target
    target_label: instance
  # Exporter IP address goes here:
  - replacement: 127.0.0.1:8088
    target_label: __address__
  scrape_interval: 5m
  static_configs:
  - labels: {}
    targets:
    # Switch admin UI URL goes here:
    - http://192.168.2.1
```


## Metrics sample

Here is a sample of provided metrics:

```
horaco_exporter_port_carrier{device="port1"} 1
horaco_exporter_port_duplex{device="port1"} 1
horaco_exporter_port_duplex_set{device="port1"} 1
horaco_exporter_port_flow_control{device="port1"} 1
horaco_exporter_port_flow_control_set{device="port1"} 1
horaco_exporter_port_receive_errs_total{device="port1"} 0
horaco_exporter_port_receive_frames_total{device="port1"} 2.0296015e+07
horaco_exporter_port_speed_bytes{device="port1"} 1.25e+08
horaco_exporter_port_speed_bytes_set{device="port1"} 0
horaco_exporter_port_state{device="port1"} 1
horaco_exporter_port_transmit_errs_total{device="port1"} 0
horaco_exporter_port_transmit_frames_total{device="port1"} 4.2407715e+07
horaco_exporter_probe_success{probe="info"} 1
horaco_exporter_probe_success{probe="port"} 1
horaco_exporter_switch_info{firmware_date="Jan 03 2024",firmware_ver="V1.9",hardware_ver="V1.1",mac="11:22:33:44:55:66",model="WAMJHJ-8125MNG"} 1
```
