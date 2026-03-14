PROJECT=mtr-exporter
VERSION?=$(shell cat VERSION)
BUILD_DATE=$(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
GIT_HASH=$(shell git rev-parse HEAD)
CONTAINER_PLATFORM?=linux/amd64

TARGETS=linux.amd64 	\
	linux.386 			\
	linux.arm64 		\
	linux.mips64 		\
	windows.amd64.exe 	\
	freebsd.amd64 		\
	darwin.amd64 		\
	darwin.arm64

BINARIES=$(addprefix bin/$(PROJECT)-$(VERSION)., $(TARGETS))
RELEASES=$(subst windows.amd64.tar.gz,windows.amd64.zip,$(foreach r,$(subst .exe,,$(TARGETS)),releases/$(PROJECT)-$(VERSION).$(r).tar.gz))

LDFLAGS:=$(LDFLAGS) -ldflags "-X main.Version=$(VERSION) -X main.BuildDate=$(BUILD_DATE) -X main.GitHash=$(GIT_HASH)"

######################################################
## release related

default: $(PROJECT)

binaries: $(BINARIES)
release: $(RELEASES)
releases: $(RELEASES)
list-releases:
	@echo $(RELEASES)|tr ' ' '\n'
clean:
	rm -f $(BINARIES) $(RELEASES)

$(PROJECT): bin/$(PROJECT)
bin/$(PROJECT): cmd/$(PROJECT) bin
	go build -v -trimpath -o $@ ./$<

bin/$(PROJECT)-$(VERSION)%:
	env GOARCH=$(subst .,,$(suffix $(subst .exe,,$@))) \
		GOOS=$(subst .,,$(suffix $(basename $(subst .exe,,$@)))) \
		CGO_ENABLED=0 \
	go build -trimpath $(LDFLAGS) -o $@ ./cmd/$(PROJECT)

releases/mtr-exporter-$(VERSION).%.zip: bin/$(PROJECT)-$(VERSION).%.exe
	mkdir -p releases
	zip -9 -j -r $@ README.md LICENSE $<
releases/$(PROJECT)-$(VERSION).%.tar.gz: bin/$(PROJECT)-$(VERSION).%
	mkdir -p releases
	tar -cf $(basename $@) README.md LICENSE && \
		tar -rf $(basename $@) --strip-components 1 $< && \
		gzip -9 $(basename $@)

bin:
	mkdir $@

CRI ?= docker
container-image:
	env DOCKER_BUILDKIT=1 $(CRI) build \
		--file Containerfile \
		--platform=$(CONTAINER_PLATFORM) \
		--build-arg VERSION=$(VERSION) \
		--tag $(CONTAINER_PLATFORM)-$(PROJECT):$(VERSION) .

######################################################
## dev related

deps-vendor:
	go mod vendor
deps-cleanup:
	go mod tidy
deps-ls:
	go list -m -mod=readonly -f '{{if not .Indirect}}{{.}}{{end}}' all
deps-ls-updates:
	go list -m -mod=readonly -f '{{if not .Indirect}}{{.}}{{end}}' -u all

test:
	go test -v ./cmd/$(PROJECT)
	go test -v ./pkg/...

generate-code:
	go generate -v ./cmd/...

compile-analysis: cmd/$(PROJECT)
	go build -gcflags '-m' ./$^

# https://github.com/nektos/act
run-github-workflow-lint:
	act -j lint --container-architecture linux/amd64
run-github-workflow-test:
	act -j test --container-architecture linux/amd64
run-github-workflow-buildLinux:
	act -j buildLinux --container-architecture linux/amd64

reports: report-golangci-lint
reports: report-vuln report-gosec report-vet

report-golangci-lint:
	@echo '####################################################################'
	golangci-lint run cmd/... pkg/...

report-vuln:
	@echo '####################################################################'
	govulncheck ./cmd/... ./pkg/...

report-grype:
	@echo '####################################################################'
	grype .

fetch-report-tools:
	go install golang.org/x/vuln/cmd/govulncheck@latest

fetch-report-tool-grype:
	go install github.com/anchore/grype@latest


.PHONY: $(PROJECT) bin/$(PROJECT) binaries releases
