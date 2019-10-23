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

if [[ -z "${ID}" ]];then
  echo "ERROR- You must provide mon id for start mon. ie: a"
  exit 1
fi

if [[ -z "${MON_SVC_ADDR}" ]];then
  echo "ERROR- You must provide mon service address for start mon. ie: 10.43.171.33"
  exit 1
fi

if [[ -z "${MON_IP}" ]];then
  echo "ERROR- You must provide mon pod ip for start mon. ie: 10.42.3.3"
  exit 1
fi

Mon_Args="--log-to-stderr=true --err-to-stderr=true --mon-cluster-log-to-stderr=true --log-stderr-prefix=debug --default-log-to-file=false --default-mon-cluster-log-to-file=false"

mon_mkfs() {
  ceph-mon --fsid=${FSID} --keyring=/etc/ceph/keyring ${Mon_Args} --mon-host=${MON_HOSTS} --mon-initial-members=${MON_MEMBERS} --id=${ID} --foreground --public-addr=${MON_SVC_ADDR} --mkfs
}

mon_run() {
  ceph-mon --fsid=${FSID} --keyring=/etc/ceph/keyring ${Mon_Args} --mon-host=${MON_HOSTS} --mon-initial-members=${MON_MEMBERS} --id=${ID} --foreground --public-addr=${MON_SVC_ADDR} --setuser-match-path=/var/lib/ceph/mon/ceph-${ID}/store.db --public-bind-addr=${MON_IP}
}

start_mon() {
  mon_mkfs
  mon_run
}

modify_conf() {
  conf="/etc/ceph/ceph.conf"
  sed -i '/fsid/d' ${conf} 
  sed -i '/mon initial members/d' ${conf} 
  sed -i '/mon host/d' ${conf}
}

modify_conf
start_mon
