run:
	go run ./cmd/server/wikiloc-layer-server.go

build:
	go build -v -x -o ./bin/wikiloc-earth-layer ./cmd/server/wikiloc-earth-layer.go

# Build instructions for docker scratch image
build-static:
	CGO_ENABLED=0 && GOOS=linux && GOARCH=amd64 && go build -a -tags netgo -ldflags '-w -extldflags "-static"' -o ./bin/wikiloc-earth-layer ./cmd/server/wikiloc-earth-layer.go

serve:
	./bin/wikiloc-earth-layer

get-icons:
	go run ./cmd/get-icons/get-icons.go

gen-kml:
	go run ./cmd/gen-kml/gen-kml.go

build-image:
	docker build -t wikiloc-earth-layer .
