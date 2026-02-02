# --- Kubemesh Build Configuration ---
DOCKER_USER  ?= shehrozdevhub
IMAGE_NAME   ?= kubemesh
VERSION      ?= v1.0.0
PLATFORMS    ?= linux/amd64,linux/arm64
BUILDER_NAME ?= kubemesh-builder

# --- Helper Variables ---
FULL_IMAGE_NAME = $(DOCKER_USER)/$(IMAGE_NAME)

.PHONY: help build-local release clean

# Default target: show help
help:
	@echo "================================================================"
	@echo "Kubemesh Build System"
	@echo "================================================================"
	@echo "Usage:"
	@echo "  make build-local  - Build for your current machine architecture"
	@echo "  make release      - Build & Push Multi-Arch (AMD64/ARM64) to Hub"
	@echo "  make clean        - Remove local build artifacts"
	@echo "================================================================"

# 1. Local Build (Fastest for testing on your current machine)
build-local:
	@echo "Building locally for current architecture..."
	docker build -t $(FULL_IMAGE_NAME):latest .
	@echo "Local build complete: $(FULL_IMAGE_NAME):latest"

# 2. Multi-Arch Release
release:
	@echo "Setting up Docker Buildx..."
	# Create and use a new builder instance if it doesn't exist
	docker buildx create --name $(BUILDER_NAME) --use || docker buildx use $(BUILDER_NAME)
	docker buildx inspect --bootstrap
	
	@echo "Starting Multi-Arch build for $(PLATFORMS)..."
	# This command builds both versions, creates the manifest, and pushes to Hub
	docker buildx build \
		--platform $(PLATFORMS) \
		-t $(FULL_IMAGE_NAME):$(VERSION) \
		-t $(FULL_IMAGE_NAME):latest \
		--push .
	
	@echo "Successfully pushed $(FULL_IMAGE_NAME):$(VERSION) to Docker Hub"
	@echo "Cleaning up builder..."
	docker buildx rm $(BUILDER_NAME)

# 3. Clean up
clean:
	@echo "Cleaning up build artifacts..."
	# Add binary paths if you build them outside of Docker
	docker rmi $(FULL_IMAGE_NAME):latest || true