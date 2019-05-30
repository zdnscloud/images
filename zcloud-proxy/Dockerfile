FROM golang:1.12.5-alpine3.9 AS build

RUN mkdir -p /go/src/github.com/zdnscloud \
    && apk add git \
    && cd /go/src/github.com/zdnscloud \
    && git clone https://github.com/zdnscloud/singlecloud.git \
    && cd singlecloud \
    && CGO_ENABLED=0 GOOS=linux go build cmd/agent/agent.go

FROM alpine:3.9.4
RUN apk --no-cache add ca-certificates
COPY --from=build /go/src/github.com/zdnscloud/singlecloud/agent /usr/local/bin/

ENTRYPOINT ["agent"]
