FROM golang:1.12.5-alpine3.9 AS build
COPY ./sleep-forever.go /go/src
RUN mkdir -p /go/src/github.com/c-bata \
    && apk add make git \
    && cd /go/src/github.com/c-bata \
    && git clone https://github.com/c-bata/kube-prompt.git \
    && cd kube-prompt \
    && git checkout v1.0.8 \
    && make build \
    && cd /go/src \
    && go build sleep-forever.go

FROM alpine:3.9.4

WORKDIR /home
ENV KUBE_VERSION="v1.13.10"
ENV PS1='zcloud-shell\$ '
COPY --from=build /go/src/github.com/c-bata/kube-prompt/kube-prompt /usr/local/bin/
COPY --from=build /go/src/sleep-forever /usr/local/bin/
RUN apk add --update -t deps curl \
    && curl -L https://storage.googleapis.com/kubernetes-release/release/${KUBE_VERSION}/bin/linux/amd64/kubectl -o /usr/local/bin/kubectl \
    && chmod +x /usr/local/bin/kubectl \
    && apk del --purge deps curl \
    && apk add bash jq \
    && rm /var/cache/apk/*

USER 1001:1001
ENTRYPOINT ["/usr/local/bin/sleep-forever"]
