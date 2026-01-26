# Stage 1: Build (Optional if you already have the binary)
# We use the official Go image to ensure a consistent environment
FROM golang:1.25.5-alpine AS builder

# Install libpcap-dev so Go can compile against the pcap headers
RUN apk add --no-cache libpcap-dev gcc musl-dev

WORKDIR /app
COPY . .

# Compile the sniffer
# We keep CGO enabled because gopacket requires it for libpcap
RUN CGO_ENABLED=1 GOOS=linux go build -o sniffer .

# Stage 2: Final Image
FROM alpine:latest

# Install libpcap (the runtime library)
RUN apk add --no-cache libpcap

# Copy the binary from the builder stage
COPY --from=builder /app/sniffer /usr/local/bin/sniffer

# Give the binary execution permissions
RUN chmod +x /usr/local/bin/sniffer

# Run the sniffer
ENTRYPOINT ["/usr/local/bin/sniffer"]