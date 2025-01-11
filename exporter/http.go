package exporter

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors/version"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
	"strings"

	"github.com/stas2k/horaco_exporter/clients"
	"github.com/stas2k/horaco_exporter/collectors"
)

type HoracoExporter bool

func (e HoracoExporter) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/metrics":
		e.serveMetrics(w, r)
	case "/":
		e.serveIndex(w, r)
	default:
		http.Error(w, "not found", http.StatusNotFound)
	}
}

func (e *HoracoExporter) serveMetrics(w http.ResponseWriter, r *http.Request) {
	target := r.URL.Query().Get("target")
	if target == "" {
		http.Error(w, "'target' parameter must be specified", http.StatusBadRequest)
		return
	}

	if !strings.Contains(target, "://") {
		target = "http://" + target
	}

	user := r.URL.Query().Get("user")
	if user == "" {
		http.Error(w, "'user' parameter must be specified", http.StatusBadRequest)
		return
	}
	password := r.URL.Query().Get("password")
	if password == "" {
		http.Error(w, "'password' parameter must be specified", http.StatusBadRequest)
		return
	}

	client := clients.NewHoracoClient(target, user, password)
	info_collector := collectors.NewSystemInfoCollector("horaco_exporter", client)
	port_collector := collectors.NewPortStatsCollector("horaco_exporter", client)

	registry := prometheus.NewRegistry()
	registry.MustRegister(version.NewCollector("horaco_exporter"))
	registry.MustRegister(info_collector)
	registry.MustRegister(port_collector)

	promhttp.HandlerFor(
		registry,
		promhttp.HandlerOpts{}).ServeHTTP(w, r)

}
func (e *HoracoExporter) serveIndex(w http.ResponseWriter, r *http.Request) {
	_, _ = w.Write([]byte(`<html>
<head><title>Horaco Exporter</title></head>
<body>
<h1>Horaco Exporter</h1>
 <form action='/metrics'>
  <label for="target">Target:</label>
	<input type="text" id="target" name="target" value="http://192.168.130.2"><br>
  <label for="user">User name:</label>
	<input type="text" id="user" name="user" value="admin"><br>
  <label for="user">Password:</label>
	<input type="text" id="password" name="password" value="admin"><br>
	<input type="submit" value="OK">
</form>
</body>
</html>
`))
}
