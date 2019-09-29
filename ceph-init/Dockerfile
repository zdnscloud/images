FROM alpine:3.10.0
LABEL maintainers="Zdns Authors"
LABEL description="K8S CSI Ceph init"

COPY shell/copy.sh /copy.sh
COPY shell/ceph_getid /ceph_getid
COPY shell/osd_volume_create.sh /osd_volume_create.sh
ENTRYPOINT ["/bin/sh"]
