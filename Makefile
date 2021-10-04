VERSION = ${shell git describe --tag}

run:
	go run ./cmd/server

serve:
	./bin/wikiloc-earth-layer

get-icons:
	go run ./cmd/get-icons

gen-kml:
	go run ./cmd/gen-kml

build:
	go build -v -x -o ./bin/wikiloc-earth-layer ./cmd/server

build-with-docker:
	docker run --mount type=bind,source="$(pwd)"/target,target=/app golang:1.17.0-bullseye make build & make gen-kml & make get-icons & wait

# Build instructions for docker scratch image
build-static:
	CGO_ENABLED=0 && GOOS=linux && GOARCH=amd64 && go build -a -tags netgo -ldflags '-w -extldflags "-static"' -o ./bin/wikiloc-earth-layer ./cmd/server

build-test-image:
	docker build --build-arg PROTOCOL=https --build-arg HOST=wikiloc-earth-layer-testing.jonamat.cloud --build-arg PORT=80 -t jonamat/wikiloc-earth-layer-testing:latest --no-cache .

push-test-image:
	docker push jonamat/wikiloc-earth-layer-testing:latest

build-prod-image:
	docker build --build-arg PROTOCOL=https --build-arg HOST=wikiloc-earth-layer.jonamat.cloud --build-arg PORT=80 -t jonamat/wikiloc-earth-layer:latest -t jonamat/wikiloc-earth-layer:${VERSION} --no-cache .

push-prod-image:
	docker push jonamat/wikiloc-earth-layer:latest
	docker push jonamat/wikiloc-earth-layer:${VERSION}

build-release:
	./scripts/release.sh
