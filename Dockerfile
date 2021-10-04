FROM golang:1.17.0-bullseye AS builder
WORKDIR /build

# Define default build args
ARG PROTOCOL=http
ARG PORT=80
ARG HOST=localhost

# Pass envs to generate static assets from builder
ENV PROTOCOL=${PROTOCOL}
ENV PORT=${PORT}
ENV HOST=${HOST}

# Envs for build
ENV CGO_ENABLED=0
ENV GOOS=linux
ENV GOARCH=amd64

# Import the codebase
COPY . .

# Create server binary, fetch & convert icons, generate init KML
RUN go run -mod vendor ./cmd/get-icons \
& go run -mod vendor ./cmd/gen-kml \
& go build -mod vendor -a -tags netgo -ldflags '-w -extldflags "-static"' -o ./bin/wikiloc-earth-layer ./cmd/server \
& wait


FROM scratch AS runner
WORKDIR /app

# Web assets from repo
COPY ./web ./web

# Server binary from builder
COPY --from=builder /build/bin/wikiloc-earth-layer ./bin/wikiloc-earth-layer

# Self-signed certificate from builder
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt

# Generated web assets from builder
COPY --from=builder /build/web/static ./web/static

# Defaults
COPY ./.env ./.env
COPY ./config.yml ./config.yml

# Run the server
ENTRYPOINT ["/app/bin/wikiloc-earth-layer"]
