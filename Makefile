VERSION=$(shell cat VERSION)
BUILD_DATE=$(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
GIT_HASH=$(shell git rev-parse HEAD)

BINARIES=bin/mtr-exporter-$(VERSION).linux.amd64 \
		 bin/mtr-exporter-$(VERSION).linux.386 \
		 bin/mtr-exporter-$(VERSION).linux.arm64 \
		 bin/mtr-exporter-$(VERSION).linux.mips64 \
		 bin/mtr-exporter-$(VERSION).windows.amd64.exe \
		 bin/mtr-exporter-$(VERSION).freebsd.amd64 \
		 bin/mtr-exporter-$(VERSION).darwin.amd64

LDFLAGS=-ldflags "-X main.Version=$(VERSION) -X main.BuildDate=$(BUILD_DATE) -X main.GitHash=$(GIT_HASH)"

mtr-exporter: cmd/mtr-exporter
	cd cmd/mtr-exporter && go build -v -o ../../$@

compile-analysis:
	cd cmd/mtr-exporter && go build -gcflags '-m'

code-quality:
	-go vet ./cmd/mtr-exporter
	-gofmt -s -d ./cmd/mtr-exporter
	-golint ./cmd/mtr-exporter
	-gocyclo ./cmd/mtr-exporter
	-ineffassign ./cmd/mtr-exporter

test:
	cd cmd/mtr-exporter && go test -v

release: $(BINARIES)

container-image:
	docker build \
		--file Dockerfile \
		--build-arg MTR_BIN=bin/mtr-exporter-$(VERSION).linux.amd64 \
		--tag mtr-exporter:$(VERSION) .

bin/mtr-exporter-$(VERSION).linux.mips64: bin
	cd cmd/mtr-exporter && env GOOS=linux GOARCH=mips64 CGO_ENABLED=0 go build $(LDFLAGS) -o ../../$@

bin/mtr-exporter-$(VERSION).linux.amd64: bin
	cd cmd/mtr-exporter && env GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build $(LDFLAGS) -o ../../$@

bin/mtr-exporter-$(VERSION).linux.386: bin
	cd cmd/mtr-exporter && env GOOS=linux GOARCH=386 CGO_ENABLED=0 go build $(LDFLAGS) -o ../../$@

bin/mtr-exporter-$(VERSION).linux.arm64: bin
	cd cmd/mtr-exporter && env GOOS=linux GOARCH=arm64 CGO_ENABLED=0 go build $(LDFLAGS) -o ../../$@

bin/mtr-exporter-$(VERSION).windows.amd64.exe: bin
	cd cmd/mtr-exporter && env GOOS=windows GOARCH=amd64 CGO_ENABLED=0 go build $(LDFLAGS) -o ../../$@

bin/mtr-exporter-$(VERSION).darwin.amd64: bin
	cd cmd/mtr-exporter && env GOOS=darwin GOARCH=amd64 CGO_ENABLED=0 go build $(LDFLAGS) -o ../../$@

bin/mtr-exporter-$(VERSION).freebsd.amd64: bin
	cd cmd/mtr-exporter && env GOOS=freebsd GOARCH=amd64 CGO_ENABLED=0 go build $(LDFLAGS) -o ../../$@

bin:
	mkdir $@
