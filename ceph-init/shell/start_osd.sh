#!/bin/bash
if [[ -z "${FSID}" ]];then
  echo "ERROR- You must provide ceph cluster fsid for start osd. ie: d6f97796-9864-4b31-b008-4a478d3b2f89"
  exit 1
fi

if [[ -z "${OSD_DEVICE}" ]];then
  echo "ERROR- You must provide a device to build your OSD ie: /dev/sdb"
  exit 1
fi

if [[ ! -e "${OSD_DEVICE}" ]]; then
  echo "ERROR- The device pointed by OSD_DEVICE ($OSD_DEVICE) doesn't exist !"
  exit 1
fi

if [[ -z "${MON_MEMBERS}" ]]; then
  echo "ERROR- You must provide mon members for start mon. ie: a,b,c"
  exit 1
fi

if [[ -z "${MON_ENDPOINT}" ]]; then
  echo "ERROR- You must provide mon endpoint for start mon. ie: a,b,c"
  exit 1
fi

if [[ -z "${ADDR}" ]]; then
  echo "ERROR- You must provide mon members for start mon. ie: a,b,c"
  exit 1
fi

if [[ -z "${CLUSTER}" ]]; then
  CLUSTER="ceph"
fi

if [[ ${OSD_BLUESTORE} -eq 1 ]]; then
  CEPH_DISK_CLI_OPTS+=(--bluestore)
fi

if [[ ! -f /var/lib/ceph/bootstrap-osd/ceph.keyring ]]; then
  mkdir -pv /var/lib/ceph/bootstrap-osd/
  mkdir -pv /var/lib/ceph/osd/
  ceph auth get client.bootstrap-osd -o /var/lib/ceph/bootstrap-osd/ceph.keyring
fi

Osd_Args="--setuser ceph --setgroup ceph --default-log-to-file false --ms-learn-addr-from-peer=false"

get_id() {
  ceph-volume lvm list --format json > /tmp/lvm.json
 # ID=$(/etc/ceph/ceph_getid /tmp/lvm.json)
  read id uuid <<< $(/etc/ceph/ceph_getid /tmp/lvm.json)
  export ID=${id}
  export UUID=${uuid}
  #echo ${ID}
}

osd_prepare() {
  #ceph-volume lvm prepare "${CEPH_DISK_CLI_OPTS[@]}" --data ${OSD_DEVICE}
  stdbuf -oL ceph-volume lvm batch --prepare "${CEPH_DISK_CLI_OPTS[@]}" --yes --osds-per-device 1 ${OSD_DEVICE}
}

osd_activate() {
  #uuid=$(cat /var/lib/ceph/osd/ceph-${ID}/fsid)
  #ceph-osd --fsid ${FSID} ${Osd_Args} --id ${ID} --cluster ${CLUSTER} --crush-location="root=default host=${OSD_NAME}" --foreground
  #ceph-osd --fsid ${FSID} ${Osd_Args} --id ${ID} --cluster ${CLUSTER}  --foreground
  stdbuf -oL ceph-volume lvm activate --no-systemd --bluestore ${ID} ${UUID}
}

osd() {
  ceph-osd --foreground --id ${ID} --fsid ${FSID} --cluster ceph --setuser ceph --setgroup ceph --default-log-to-file false --ms-learn-addr-from-peer=false --crush-location="root=default host=${OSD_NAME}" 
}

conf() {
cat > /etc/ceph/ceph.conf << EOF
[global]
fsid                      = ${FSID}
mon initial members       = ${MON_MEMBERS}
mon host                  = ${MON_ENDPOINT}
public addr               = ${ADDR}
cluster addr              = ${ADDR}
mon keyvaluedb            = rocksdb
mon_allow_pool_delete     = true
mon_max_pg_per_osd        = 1000
debug default             = 0
debug rados               = 0
debug mon                 = 0
debug osd                 = 0
debug bluestore           = 0
debug filestore           = 0
debug journal             = 0
debug leveldb             = 0
filestore_omap_backend    = rocksdb
osd pg bits               = 11
osd pgp bits              = 11
osd pool default size     = 1
osd pool default min size = 1
osd pool default pg num   = 100
osd pool default pgp num  = 100
osd objectstore           = filestore
crush location            = root=default host=${OSD_NAME}
rbd_default_features      = 3
fatal signal handlers     = false

[osd.0]
keyring              = /var/lib/ceph/osd/ceph-${ID}/keyring
bluestore block path = /var/lib/ceph/osd/ceph-${ID}/block
EOF
}

osd_start() {
  get_id
  if [ -n "${ID}" -a -n "${UUID}" ];then
    osd_activate
    conf
    osd
  fi
}

osd_start
