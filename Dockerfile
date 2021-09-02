FROM golang:1.17.0-bullseye AS builder
WORKDIR /build
COPY . .
# Create server binary, fetch & convert icons, generate init KML
RUN make build-static & make get-icons & make gen-kml & wait

FROM scratch AS runner
WORKDIR /app
# Server binary from builder
COPY --from=builder /build/bin/wikiloc-earth-layer ./bin/wikiloc-earth-layer
# Self-signed certificate from builder
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
# Static & templates from repo
COPY ./web ./web
# Def envs
COPY ./.env .

ENTRYPOINT ["/app/bin/wikiloc-earth-layer"]
