##
## -- runtime environment
##

FROM    golang:1.17.5-alpine3.15 AS build-env

ARG     MTR_BIN=bin/mtr-exporter-0.1.0.linux.amd64

ADD     . /src/mtr-exporter
RUN     apk add -U --no-cache make git
RUN     make -C /src/mtr-exporter $MTR_BIN

##
## -- runtime environment
##

FROM    alpine:3.15 AS rt-env

ARG     MTR_BIN=bin/mtr-exporter-0.1.0.linux.amd64
RUN     apk add -U --no-cache mtr
COPY    --from=build-env /src/mtr-exporter/$MTR_BIN /usr/local/bin/mtr-exporter

EXPOSE  8080
ENTRYPOINT ["/usr/local/bin/mtr-exporter"]
