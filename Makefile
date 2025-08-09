BINARY_NAME=unfurl
BINARY_PATH=./bin/$(BINARY_NAME)
IMAGE=mtsfy/unfurl
TAG=latest

# Build for my current OS (macOS local-dev)
build:
	mkdir -p ./bin
	CGO_ENABLED=0 GOOS=darwin go build -a -installsuffix cgo -o $(BINARY_PATH) ./cmd/main.go

# Build specifically for Linux (docker/prod)
build-linux:
	mkdir -p ./bin
	CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o $(BINARY_PATH) ./cmd/main.go

# Run the application locally 
run: build
	$(BINARY_PATH)

# Build Docker image from Dockerfile 
docker-build:
	docker build -t $(IMAGE):$(TAG) .

# Run Docker container on port 8080
docker-run:
	docker run -p 8080:8080 $(IMAGE):$(TAG)

# Push to Docker Hub 
docker-push:
	docker push $(IMAGE):$(TAG)

# Tag and push with version (VERSION=v1.0.0 make docker-push-version)
docker-push-version:
	docker tag $(IMAGE):$(TAG) $(IMAGE):$(VERSION)
	docker push $(IMAGE):$(VERSION)

# Remove all build artifacts and clean bin directory
clean:
	rm -rf ./bin 

.PHONY: build build-linux run docker-build docker-run docker-push docker-push-version clean