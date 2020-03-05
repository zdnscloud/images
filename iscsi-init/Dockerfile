FROM ubuntu:14.04
LABEL maintainers="Zdns Authors"
LABEL description="K8S CSI Ceph init"

RUN apt-get update && apt-get install -y open-iscsi lvm2 lsscsi -y
COPY shell/init.sh /init.sh
ENTRYPOINT ["/bin/sh"]
