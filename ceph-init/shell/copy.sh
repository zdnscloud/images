#!/bin/bash
cp /host/ceph/ceph.conf /host/etc/ceph
cp /host/ceph/ceph.client.admin.keyring /host/etc/ceph
cp /host/ceph/ceph.mon.keyring /host/etc/ceph
cp /host/ceph/osd_volume_prepare.sh /host/etc/ceph
cp /ceph_getid /host/etc/ceph
cp /osd_volume_create.sh /host/etc/ceph
