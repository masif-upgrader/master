FROM golang:alpine as build
RUN ["apk", "add", "upx"]

ADD . /master
WORKDIR /master
ENV CGO_ENABLED 0

RUN ["go", "mod", "download"]
RUN ["go", "generate", "./..."]
RUN ["go", "build", "-ldflags", "-s -w", "."]
RUN ["upx", "master"]

WORKDIR /master/docker
RUN ["go", "build", "-ldflags", "-s -w", "."]
RUN ["upx", "docker"]


FROM scratch
COPY --from=build /master/master /master
COPY --from=build /master/docker/docker /entrypoint
CMD ["/entrypoint"]
