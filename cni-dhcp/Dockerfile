FROM alpine:3.10.0
LABEL maintainers="Zdns Authors"
LABEL description="K8S CNI MacVlan+DHCP"

RUN mkdir -p /opt/cni/bin
WORKDIR /opt/cni/bin
RUN wget https://github.com/containernetworking/plugins/releases/download/v0.7.6/cni-plugins-amd64-v0.7.6.tgz
RUN tar fvxz cni-plugins-amd64-v0.7.6.tgz
RUN rm cni-plugins-amd64-v0.7.6.tgz
COPY shell/install-cni.sh /install-cni.sh
ENTRYPOINT ["/bin/sh"]
