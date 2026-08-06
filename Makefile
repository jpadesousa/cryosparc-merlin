BINARY := cryosparc-merlin 

.PHONY: all fmt vet lint test build clean

all: fmt vet lint test build

fmt:
	go fmt ./...

vet:
	go vet ./...

lint:
	/data/user/agosti_j/go/bin/golangci-lint run

test:
	go test ./...

build:
	mkdir -p bin
	go build -o bin/$(BINARY)

clean:
	rm -rf bin/
