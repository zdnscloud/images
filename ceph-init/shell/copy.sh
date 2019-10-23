#!/bin/bash
cp /host/ceph/ceph.conf /host/etc/ceph
cp /host/ceph/ceph.client.admin.keyring /host/etc/ceph
cp /host/ceph/keyring /host/etc/ceph
cp /ceph_getid /host/etc/ceph
cp /start_mon.sh /host/etc/ceph
cp /start_osd.sh /host/etc/ceph
cp /start_mgr.sh /host/etc/ceph
cp /start_mds.sh /host/etc/ceph
cp /prepare.sh /host/etc/ceph
