FROM golang:1.11 as build

RUN go get github.com/golang/dep \
	&& go install github.com/golang/dep/...

ADD . /go/src/github.com/masif-upgrader/master

RUN cd /go/src/github.com/masif-upgrader/master \
	&& /go/bin/dep ensure \
	&& go generate \
	&& go install .

FROM debian:9

COPY --from=build /go/bin/master /usr/local/bin/masif-upgrader-master
COPY --from=ochinchina/supervisord:latest /usr/local/bin/supervisord /usr/local/bin/

COPY --from=masifupgrader/common /pki-master/keys /pki-master
COPY --from=masifupgrader/common /pki-agent/keys /pki-agent
COPY _docker/config.ini /etc/masif-upgrader-master.ini
COPY _docker/supervisord.conf /etc/

CMD ["/usr/local/bin/supervisord", "-c", "/etc/supervisord.conf"]
