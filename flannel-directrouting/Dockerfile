FROM golang:1.12.5-alpine3.9 AS build
  
RUN mkdir -p /go/src/github.com/zdnscloud/images/flannel-directrouting
COPY . /go/src/github.com/zdnscloud/images/flannel-directrouting

WORKDIR /go/src/github.com/zdnscloud/images/flannel-directrouting
RUN CGO_ENABLED=0 GOOS=linux go build -o flanneld

FROM alpine:3.10.0

LABEL maintainers="Zdns Authors"
LABEL description="flannel vxlan directrouting sidecar for route check"
RUN apk update
RUN apk add --no-cache iproute2 net-tools ca-certificates iptables && update-ca-certificates
RUN apk add wireguard-tools --no-cache --repository http://dl-cdn.alpinelinux.org/alpine/edge/testing
COPY --from=build /go/src/github.com/zdnscloud/images/flannel-directrouting/flanneld /flanneld
ENTRYPOINT ["/bin/sh"]
