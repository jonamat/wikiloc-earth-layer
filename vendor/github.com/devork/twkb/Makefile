TAGS := "${GOOS} amd64"
BENCHTIME := "1m"

clean:
	go clean

deps:
	go get github.com/stretchr/testify/assert

build:
	go clean
	go build -tags ${TAGS}

test:
	go test -v
