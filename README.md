# wargame metric exporter for Promehteus

**Wargame** is a real-time strategy game service. This exporter collects custom metrics from a server.

It reads JSON snapshots created by the server and exposes them.

<br>

Even if your server is not C++, this exporter works for any service that creates structured metric snapshots.

<br>

### Quick Start

```
# Default: Port 9100, Metrics file: /mnt/metrics/metrics_snapshot.json
./metrics-exporter

# Using environment variables
EXPORTER_PORT=9200 METRICS_FILE=/data/metrics.json ./metrics-exporter
```

<br>

## Usage in kubernetes

The `/docs` directory contains examples for different deployment patterns.

Select the one that fits your service scale:

1. **Sidecar Pattern** – Run the exporter inside the same Pod as the game server.

2. **Shared PVC Pattern** – Run the exporter in a separate Pod and share a PVC with the game server.

<br>

## How to Configure the server? (C++20)

We want to **minimize lock issues** and create snapshots safely.

Use a global metrics structure and a timer thread to write JSON snapshots using `nlohmann::json`.

<br>

```
#include <thread>
#include <chrono>
#include <fstream>
#include <nlohmann/json.hpp>
#include <atomic>

struct Metrics { // custom
    std::atomic<int> activePlayers{0};
    ...
};

void DumpMetricsSnapshot(const std::string& path) {
    nlohmann::json j;
    ...
}

void MetricsSnapshotThread() {
    while (true) {
        DumpMetricsSnapshot("/mnt/metrics/metrics_snapshot.json");
        std::this_thread::sleep_for(std::chrono::seconds(1));
    }
}

// server init:
std::thread metricsThread(MetricsSnapshotThread);
metricsThread.detach();
```

<br>

## Makefile Commands
```
# Build the local Go binary
make build

# Build the Docker image
make docker-build

# Run the container (default port/file)
make run

# Change port or metrics file
make run EXPORTER_PORT=9200 METRICS_FILE=/data/metrics.json

# View logs
make logs

# Clean up
make dist-clean
```
