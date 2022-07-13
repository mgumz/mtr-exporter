VERSION=$(shell cat VERSION)
BUILD_DATE=$(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
GIT_HASH=$(shell git rev-parse HEAD)

BUILDS=linux.amd64 linux.386 linux.arm64 linux.mips64 windows.amd64.exe freebsd.amd64 darwin.amd64 darwin.arm64
BINARIES=$(addprefix bin/mtr-exporter-$(VERSION)., $(BUILDS))

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

bin/mtr-exporter-$(VERSION).linux.%: bin
	cd cmd/mtr-exporter && env GOOS=linux GOARCH=$* CGO_ENABLED=0 go build $(LDFLAGS) -o ../../$@

bin/mtr-exporter-$(VERSION).darwin.%: bin
	cd cmd/mtr-exporter && env GOOS=darwin GOARCH=$* CGO_ENABLED=0 go build $(LDFLAGS) -o ../../$@

bin/mtr-exporter-$(VERSION).windows.%.exe: bin
	cd cmd/mtr-exporter && env GOOS=windows GOARCH=$* CGO_ENABLED=0 go build $(LDFLAGS) -o ../../$@

bin/mtr-exporter-$(VERSION).freebsd.%: bin
	cd cmd/mtr-exporter && env GOOS=freebsd GOARCH=$* CGO_ENABLED=0 go build $(LDFLAGS) -o ../../$@

bin:
	mkdir $@
