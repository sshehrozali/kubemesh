# --- Stage 1: Build ---
# We use BUILDPLATFORM to ensure the builder itself runs on the host's native speed
FROM --platform=$BUILDPLATFORM golang:1.25-alpine AS builder

# 1. Install 'xx' - the magic bridge for cross-compiling CGO
COPY --from=tonistiigi/xx / /

# 2. Install native build tools for the BUILDER (the machine running the build)
# We use clang because it is a native cross-compiler, unlike gcc
RUN apk add --no-cache clang lld libpcap-dev gcc musl-dev

# 3. Setup the TARGET environment
ARG TARGETPLATFORM
# This tells 'xx' to pull the libpcap and musl headers for the target (e.g., AMD64)
RUN xx-apk add --no-cache libpcap-dev musl-dev

WORKDIR /app

# Cache Go modules
COPY go.mod go.sum ./
RUN go mod download

COPY . .

# 4. Build with xx-go
# xx-go automatically sets CC, CXX, and other env vars for the target
ARG TARGETOS TARGETARCH
RUN CGO_ENABLED=1 xx-go build \
    -ldflags="-s -w" \
    -o /kubemesh ./cmd/kubemesh/main.go && \
    xx-verify /kubemesh

# --- Stage 2: Final Runtime ---
# This stage defaults to the TARGETPLATFORM automatically
FROM alpine:3.19

# Install only the runtime library (no -dev headers needed here)
RUN apk add --no-cache libpcap ca-certificates

# Copy the optimized binary
COPY --from=builder /kubemesh /usr/local/bin/kubemesh

# Product Metadata
LABEL org.opencontainers.image.title="KubeMesh"
LABEL org.opencontainers.image.description="Zero-sidecar K8s traffic observability"

# Default Production Config
ENV TRAFFIC_PORT=80

ENTRYPOINT ["/usr/local/bin/kubemesh"]