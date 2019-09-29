#!/bin/bash
set -e
if [[ -z "${OSD_DEVICE}" ]];then
  log "ERROR- You must provide a device to build your OSD ie: /dev/sdb"
  exit 1
fi

if [[ ! -e "${OSD_DEVICE}" ]]; then
  log "ERROR- The device pointed by OSD_DEVICE ($OSD_DEVICE) doesn't exist !"
  exit 1
fi
if [[ ! -f /var/lib/ceph/bootstrap-osd/ceph.keyring ]]; then
  ceph auth get client.bootstrap-osd -o /var/lib/ceph/bootstrap-osd/ceph.keyring
fi
if [[ ${OSD_BLUESTORE} -eq 1 ]]; then
   CEPH_DISK_CLI_OPTS+=(--bluestore)
fi
ceph-volume lvm list --format json > /tmp/lvm.json
ID=$(/etc/ceph/ceph_getid /tmp/lvm.json)
if [[ -z "${ID}" ]];then
  ceph-volume lvm prepare "${CEPH_DISK_CLI_OPTS[@]}" --data ${OSD_DEVICE}
  ceph-volume lvm list --format json > /tmp/lvm.json
  ID=$(/etc/ceph/ceph_getid /tmp/lvm.json)
fi
if [[ -n "${ID}" ]];then
  export OSD_ID=${ID}
  /opt/ceph-container/bin/entrypoint.sh osd_ceph_volume_activate
fi
