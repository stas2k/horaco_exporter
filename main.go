package main

import (
	"flag"
	"log"
	"net/http"
	"strconv"

	"github.com/stas2k/horaco_exporter/exporter"
)

func main() {
	var port = flag.Int("port", 8088, "port to listen for Prometheus scrapes on")
	flag.Parse()

	e := exporter.HoracoExporter(true)
	http.Handle("/", e)

	log.Printf("Beginning to serve on port :" + strconv.Itoa(*port))
	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(*port), nil))
}
