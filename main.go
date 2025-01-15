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
	var hosts = flag.String("hosts", "", "space delimited file containing target credentials")
	flag.Parse()

	auth := &exporter.AuthEndpoints{}
	var err error
	if *hosts != "" {
		auth, err = exporter.ParseAuthFile(*hosts)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		log.Println("WARNING! Running in insecure mode!")
		log.Println("Please use the -hosts option to limit the hosts reachable by this exporter")
	}

	e := exporter.HoracoExporter(*auth)
	http.Handle("/", e)

	log.Printf("Beginning to serve on port :" + strconv.Itoa(*port))
	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(*port), nil))
}
