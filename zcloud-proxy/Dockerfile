FROM golang:1.13.7-alpine3.11 AS build

ARG goproxy
ENV GOPROXY=$goproxy

RUN mkdir -p /go/src/github.com/zdnscloud/images/zcloud-proxy
COPY . /go/src/github.com/zdnscloud/images/zcloud-proxy
RUN cd /go/src/github.com/zdnscloud/images/zcloud-proxy \
    && CGO_ENABLED=0 GOOS=linux go build agent.go

FROM alpine:3.11.5
RUN apk --no-cache add ca-certificates
COPY --from=build /go/src/github.com/zdnscloud/images/zcloud-proxy/agent /usr/local/bin/
ENTRYPOINT ["agent"]
