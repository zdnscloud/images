#!/bin/bash
set -e

if [[ -z "${FSID}" ]];then
  echo "ERROR- You must provide ceph cluster fsid for start mon. ie: d6f97796-9864-4b31-b008-4a478d3b2f89"
  exit 1
fi

if [[ -z "${MON_HOSTS}" ]];then
  echo "ERROR- You must provide mon hosts for start mon. ie: [v2:10.43.171.33:3300,v1:10.43.171.33:6789],[v2:10.43.16.94:3300,v1:10.43.16.94:6789]"
  exit 1
fi

if [[ -z "${MON_MEMBERS}" ]];then
  echo "ERROR- You must provide mon members for start mon. ie: a,b,c"
  exit 1
fi

if [[ -z "${MDS_NAME}" ]];then
  echo "ERROR- You must provide mon pod ip for start mon. ie: ceph-mds-a"
  exit 1
fi

if [[ -z "${REPLICATION}" ]];then
  export REPLICATION=1
fi

Mds_Args="--log-to-stderr=true --err-to-stderr=true --mon-cluster-log-to-stderr=true --log-stderr-prefix=debug --default-log-to-file=false --default-mon-cluster-log-to-file=false"
Key="/etc/ceph/keyring"

create_key() {
  ceph auth get-or-create mds.${MDS_NAME} osd 'allow rwx' mds 'allow' mon 'allow profile mds' -o ${Key}
}

create_pool() {
  if [ "$CEPHFS_CREATE" -eq 1 ]; then
    ceph osd pool create "${CEPHFS_DATA_POOL}" "${CEPHFS_DATA_POOL_PG}"
    ceph osd pool set ${CEPHFS_DATA_POOL} size ${REPLICATION}
    ceph osd pool create "${CEPHFS_METADATA_POOL}" "${CEPHFS_METADATA_POOL_PG}"
    ceph osd pool set ${CEPHFS_METADATA_POOL} size ${REPLICATION}
    ceph fs new "${CEPHFS_NAME}" "${CEPHFS_METADATA_POOL}" "${CEPHFS_DATA_POOL}"
  fi
}

mds_run() {
  ceph-mds --fsid=${FSID} --keyring=${Key} ${Mds_Args} --mon-host=${MON_HOSTS} --mon-initial-members=${MON_MEMBERS} --id=${MDS_NAME} --setuser=ceph --setgroup=ceph --foreground 
}

start_mds() {
  create_key
  create_pool
  mds_run
}

start_mds
