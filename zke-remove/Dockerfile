FROM golang:1.12.5-alpine3.9 AS build

RUN mkdir -p /go
COPY . /go

WORKDIR /go
RUN CGO_ENABLED=0 GOOS=linux go build remove.go

FROM alpine:3.10.0

LABEL maintainers="Kubernetes Authors"
LABEL description="Cluster Remove"

RUN apk update && apk --no-cache add util-linux blkid file util-linux e2fsprogs lvm2 udev sgdisk device-mapper

COPY --from=build /go/remove /remove
ENTRYPOINT ["/remove"]
