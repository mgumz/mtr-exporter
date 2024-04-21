##
## -- runtime environment
##

FROM    golang:1.22.2-alpine3.19 AS build-env

#       https://github.com/docker-library/official-images#multiple-architectures
#       https://docs.docker.com/engine/reference/builder/#automatic-platform-args-in-the-global-scope
ARG     TARGETPLATFORM
ARG     TARGETOS
ARG     TARGETARCH

ARG     VERSION=latest

ADD     . /src/mtr-exporter
RUN     apk add -U --no-cache make git
RUN     make LDFLAGS="-ldflags -w" -C /src/mtr-exporter bin/mtr-exporter-$VERSION.$TARGETOS.$TARGETARCH

##
## -- runtime environment
##

FROM    alpine:3.19 AS rt-env

RUN     apk add -U --no-cache mtr tini
COPY    --from=build-env /src/mtr-exporter/bin/* /usr/bin/mtr-exporter

EXPOSE  8080
ENTRYPOINT ["/sbin/tini", "--"]
CMD     ["/usr/bin/mtr-exporter"]
