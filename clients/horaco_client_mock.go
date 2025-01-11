package clients

import (
	"fmt"
	"net/http"
	"net/url"
)

func mockMainHandler(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/login.cgi":
		mockLoginHandler(w, r)
	case "/port.cgi":
		mockPortHandler(w, r)
	case "/info.cgi":
		mockInfoHandler(w, r)
	default:
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "not found")
	}
}

func mockLoginHandler(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		return
}

func mockPortHandler(w http.ResponseWriter, r *http.Request) {
	query, err := url.ParseQuery(r.URL.RawQuery)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "invalid request")
		return
	}
	page := query.Get("page")

	if page == "stats" {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(PORT_STAT_RES))
	} else if page == "" {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(PORT_RES))
	} else {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "invalid request")
		return
	}
}

func mockInfoHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(INFO_RES))
}
