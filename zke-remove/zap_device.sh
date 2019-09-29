#!/bin/bash
pvs=$(pvs --noheadings -o pv_name,vg_name|awk '{if ($2 == "k8s"||$2 ~ "ceph-") print $1}') 
vgs=$(vgs --noheadings -o vg_name|awk '{print $1}' |grep -E "^k8s$|^ceph-")
for vg in ${vgs};do
	vgremove -v -f ${vg}
done
for pv in ${pvs};do
	pvremove -v -y ${pv}
	wipefs --all ${pv}
done
