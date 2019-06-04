FROM golang:1.12.5-alpine3.9 AS build

RUN mkdir -p /go/src/github.com/zdnscloud/images/zke-tools
COPY ./ /go/src/github.com/zdnscloud/images/zke-tools
WORKDIR /go/src/github.com/zdnscloud/images/zke-tools
RUN CGO_ENABLED=0 GOOS=linux go build main.go \
    && mv main etcd-backup


FROM nginx:1.14.0-alpine

RUN apk -U --no-cache add bash curl wget ca-certificates tar sysstat\
    && mkdir -p /opt/zke-tools/bin /etc/confd \
    && curl -sLf https://github.com/kelseyhightower/confd/releases/download/v0.16.0/confd-0.16.0-linux-amd64 > /usr/bin/confd \
    && chmod +x /usr/bin/confd \
    && curl -sLf https://get.docker.com/builds/Linux/x86_64/docker-1.12.3.tgz | tar xvzf - -C /opt/zke-tools/bin --strip-components=1 docker/docker \
    && chmod +x /opt/zke-tools/bin/docker \
    && apk del curl

RUN mkdir -p /opt/cni/bin
RUN wget -q -O - https://github.com/containernetworking/cni/releases/download/v0.4.0/cni-amd64-v0.4.0.tgz | tar xzf - -C /tmp
RUN wget -q -O /tmp/portmap https://github.com/projectcalico/cni-plugin/releases/download/v1.9.1/portmap

RUN wget -q -O - https://github.com/coreos/etcd/releases/download/v3.0.17/etcd-v3.0.17-linux-amd64.tar.gz | tar xzf - -C /tmp && \
    mv /tmp/etcd-*/etcdctl /usr/local/bin/etcdctl && \
    rm -rf /tmp/etcd-* && rm -f /etcd-*.tar.gz && \
    apk del wget && \
    mkdir  /etc/docker && \
    mkdir  /docker-certs-tmp

COPY templates /etc/confd/templates/
COPY conf.d /etc/confd/conf.d/
COPY cert-deployer nginx-proxy registry-cert /usr/bin/
COPY entrypoint.sh /opt/zke-tools/
COPY --from=build /go/src/github.com/zdnscloud/images/zke-tools/etcd-backup /opt/zke-tools

VOLUME /opt/zke-tools
CMD ["/bin/bash"]
