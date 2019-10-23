#!/bin/bash
mkdir -pv /var/lib/ceph/mgr/ceph-${MGR_NAME}/
ceph auth get-or-create mgr.${MGR_NAME} mon 'allow profile mgr' osd 'allow *' mds 'allow *' -o /var/lib/ceph/mgr/ceph-${MGR_NAME}/keyring
/usr/bin/ceph-mgr --cluster ceph --default-log-to-file=false --default-mon-cluster-log-to-file=false --setuser ceph --setgroup ceph -d -i ${MGR_NAME}
