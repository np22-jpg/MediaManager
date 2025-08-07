FROM docker.io/library/golang:1.24-alpine AS metadata_relay-build

RUN apk add --no-cache git ca-certificates tzdata

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download
COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags="-w -s" -trimpath -o metadata_relay .


FROM docker.io/library/busybox:stable-musl AS runtime

COPY --from=metadata_relay-build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=metadata_relay-build /usr/share/zoneinfo /usr/share/zoneinfo

RUN adduser -D -u 1000 -g 1000 appuser

COPY --from=metadata_relay-build /app/metadata_relay /app/metadata_relay

# Set ownership and make executable
RUN chmod +x /app/metadata_relay

# Switch to non-root user
USER 1000:1000

# Expose port
EXPOSE 8080

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8080/ || exit 1

# Run the application
CMD ["/app/metadata_relay"]
