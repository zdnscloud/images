package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

func main() {
	var id string
	filePth := os.Args[1]
	f, err := os.Open(filePth)
	if err != nil {
		fmt.Println(id)
	}
	info, _ := ioutil.ReadAll(f)
	res := make(map[string][]lvm)
	json.Unmarshal(info, &res)
	for _, v := range res {
		for _, d := range v[0].Dev {
			if d != os.Getenv("OSD_DEVICE") {
				continue
			}
			id = v[0].Tags.OSD_ID
		}
	}
	fmt.Println(id)
}

type lvm struct {
	Dev     []string `json:"devices"`
	LV_NAME string   `json:"lv_name"`
	LV_PATH string   `json:"lv_path"`
	LV_SIZE string   `json:"lv_size"`
	LV_TAGS string   `json:"lv_tags"`
	LV_UUID string   `json:"lv_uuid"`
	Name    string   `json:"name"`
	Path    string   `json:"path"`
	Tags    tags     `json:"tags"`
	Type    string   `json:"type"`
	VG_NAME string   `json:"vg_name"`
}

type tags struct {
	BLOCK_Dev     string `json:"ceph.block_device"`
	BLOCK_UUID    string `json:"ceph.block_uuid"`
	SECRET        string `json:"ceph.cephx_lockbox_secret"`
	CLUSTER_FSID  string `json:"ceph.cluster_fsid"`
	CLUSTER_NAME  string `json:"ceph.cluster_name"`
	CLUSTER_CLASS string `json:"ceph.crush_device_class"`
	ENCRY         string `json:"ceph.encrypted"`
	OSD_FSID      string `json:"ceph.osd_fsid"`
	OSD_ID        string `json:"ceph.osd_id"`
	TYPE          string `json:"ceph.type"`
	VDO           string `json:"ceph.vdo"`
}
