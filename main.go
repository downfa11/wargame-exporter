package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

type Metrics struct { // custom metrics
	ActivePlayers  int `json:"active_players"`
	MatchQueueSize int `json:"match_queue_size"`
	AvgLatencyMs   int `json:"avg_latency_ms"`
}

func metricsHandler(metricsFile string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data, err := ioutil.ReadFile(metricsFile)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		var m Metrics
		if err := json.Unmarshal(data, &m); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "text/plain; version=0.0.4")
		fmt.Fprintf(w, "active_players %d\n", m.ActivePlayers)
		fmt.Fprintf(w, "match_queue_size %d\n", m.MatchQueueSize)
		fmt.Fprintf(w, "avg_latency_ms %d\n", m.AvgLatencyMs)
	}
}

func main() {
	port := os.Getenv("EXPORTER_PORT")
	if port == "" {
		port = "9100"
	}

	metricsFile := os.Getenv("METRICS_FILE")
	if metricsFile == "" {
		metricsFile = "/mnt/metrics/metrics_snapshot.json"
	}

	http.HandleFunc("/metrics", metricsHandler(metricsFile))
	addr := ":" + port
	log.Printf("Exporter listening on %s (metrics file: %s)\n", addr, metricsFile)
	log.Fatal(http.ListenAndServe(addr, nil))
}
