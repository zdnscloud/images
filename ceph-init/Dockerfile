FROM alpine:3.10.0
LABEL maintainers="Zdns Authors"
LABEL description="K8S CSI Ceph init"

COPY shell/copy.sh /copy.sh
COPY shell/ceph_getid /ceph_getid
COPY shell/start_mon.sh /start_mon.sh
COPY shell/start_osd.sh /start_osd.sh
COPY shell/start_mgr.sh /start_mgr.sh
COPY shell/start_mds.sh /start_mds.sh
COPY shell/prepare.sh /prepare.sh
ENTRYPOINT ["/bin/sh"]
