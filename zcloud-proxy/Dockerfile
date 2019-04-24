FROM golang:alpine AS build

RUN mkdir -p /go/src/github.com/zdnscloud \
    && apk add git \
    && cd /go/src/github.com/zdnscloud \
    && git clone https://github.com/zdnscloud/singlecloud.git \
    && cd singlecloud \
    && CGO_ENABLED=0 GOOS=linux go build cmd/agent/agent.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates
COPY --from=build /go/src/github.com/zdnscloud/singlecloud/agent /usr/local/bin/

ENTRYPOINT ["agent"]
