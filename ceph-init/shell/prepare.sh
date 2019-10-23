#!/bin/bash
ceph auth get client.bootstrap-osd -o /var/lib/ceph/bootstrap-osd/ceph.keyring
for dev in ${OSD_DEVICES};do
  stdbuf -oL ceph-volume lvm batch --prepare --bluestore --yes --osds-per-device 1 ${dev}
done
