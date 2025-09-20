package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

// Metrics – custom metrics snapshot
type Metrics struct {
	ActivePlayers          int     `json:"active_players"`            // 현재 사용자 수
	ActiveRooms            int     `json:"active_rooms"`              // 현재 진행 중인 방 수
	CPUUsagePercent        float64 `json:"cpu_usage_percent"`         // CPU 사용률 (%)
	MemoryUsageBytes       int64   `json:"memory_usage_bytes"`        // RAM 사용량 (bytes)
	AvgLatencyMs           int     `json:"avg_latency_ms"`            // 평균 네트워크 지연 (ms)
	PacketLossRatioPercent float64 `json:"packet_loss_ratio_percent"` // 패킷 손실률 (%)
	KafkaMatchingMessages  int     `json:"kafka_matching_messages"`   // 매칭 데이터 수신 메시지 수(누적)
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

		m.MemoryUsageBytes = int64(m.MemoryUsageBytes) * 1024 * 1024
		w.Header().Set("Content-Type", "text/plain; version=0.0.4")
		w.Write([]byte(
			fmt.Sprintf("active_players %d\n", m.ActivePlayers) +
				fmt.Sprintf("active_rooms %d\n", m.ActiveRooms) +
				fmt.Sprintf("cpu_usage_percent %.2f\n", m.CPUUsagePercent) +
				fmt.Sprintf("memory_usage_bytes %d\n", m.MemoryUsageBytes) +
				fmt.Sprintf("avg_latency_ms %d\n", m.AvgLatencyMs) +
				fmt.Sprintf("packet_loss_ratio_percent %.2f\n", m.PacketLossRatioPercent) +
				fmt.Sprintf("kafka_matching_messages %d\n", m.KafkaMatchingMessages),
		))
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
