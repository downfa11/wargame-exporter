IMAGE_NAME=iocp-metrics-exporter
TAG=latest

BINARY_NAME=metrics-exporter

EXPORTER_PORT?=9100
METRICS_FILE?=/mnt/metrics/metrics_snapshot.json

build:
	@echo ">> Building Go binary..."
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o $(BINARY_NAME) main.go

docker-build:
	@echo ">> Building Docker image..."
	docker build -t $(IMAGE_NAME):$(TAG) .

run:
	@echo ">> Running container..."
	docker run -d \
		-p $(EXPORTER_PORT):$(EXPORTER_PORT) \
		-e EXPORTER_PORT=$(EXPORTER_PORT) \
		-e METRICS_FILE=$(METRICS_FILE) \
		-v $$(pwd)/metrics:/mnt/metrics \
		--name $(IMAGE_NAME) \
		$(IMAGE_NAME):$(TAG)

logs:
	docker logs -f $(IMAGE_NAME)

clean:
	@echo ">> Stopping and removing container..."
	-docker stop $(IMAGE_NAME)
	-docker rm $(IMAGE_NAME)

dist-clean: clean
	@echo ">> Removing binary..."
	-rm -f $(BINARY_NAME)