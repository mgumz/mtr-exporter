VERSION=$(shell cat VERSION)
BUILD_DATE=$(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
GIT_HASH=$(shell git rev-parse HEAD)
CONTAINER_PLATFORM?=linux/amd64

BUILDS=linux.amd64 linux.386 linux.arm64 linux.mips64 windows.amd64.exe freebsd.amd64 darwin.amd64 darwin.arm64
BINARIES=$(addprefix bin/mtr-exporter-$(VERSION)., $(BUILDS))

LDFLAGS=-ldflags "-X main.Version=$(VERSION) -X main.BuildDate=$(BUILD_DATE) -X main.GitHash=$(GIT_HASH)"

mtr-exporter: cmd/mtr-exporter
	go build -v -o $@ ./$^

######################################################
## release related

release: $(BINARIES)

bin/mtr-exporter-$(VERSION).linux.%: bin
	env GOOS=linux GOARCH=$* CGO_ENABLED=0 go build $(LDFLAGS) -o $@ ./cmd/mtr-exporter

bin/mtr-exporter-$(VERSION).darwin.%: bin
	env GOOS=darwin GOARCH=$* CGO_ENABLED=0 go build $(LDFLAGS) -o $@ ./cmd/mtr-exporter

bin/mtr-exporter-$(VERSION).windows.%.exe: bin
	env GOOS=windows GOARCH=$* CGO_ENABLED=0 go build $(LDFLAGS) -o $@ ./cmd/mtr-exporter

bin/mtr-exporter-$(VERSION).freebsd.%: bin
	env GOOS=freebsd GOARCH=$* CGO_ENABLED=0 go build $(LDFLAGS) -o $@ ./cmd/mtr-exporter

bin:
	mkdir $@


container-image:
	env DOCKER_BUILDKIT=1 docker build \
		--file Dockerfile \
		--platform=$(CONTAINER_PLATFORM) \
		--build-arg VERSION=$(VERSION) \
		--tag $(CONTAINER_PLATFORM)-mtr-exporter:$(VERSION) .

######################################################
## dev related

compile-analysis: cmd/mtr-exporter
	go build -gcflags '-m' ./$^

code-quality:
	-go vet ./cmd/mtr-exporter
	-gofmt -s -d ./cmd/mtr-exporter
	-golint ./cmd/mtr-exporter
	-gocyclo ./cmd/mtr-exporter
	-ineffassign ./cmd/mtr-exporter

test:
	go test -v ./cmd/mtr-exporter
