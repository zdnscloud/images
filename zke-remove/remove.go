package main

import (
	"bytes"
	"fmt"
	"net"
	"os"
	"os/exec"
	"strings"
)

func main() {
	UnmountPodVolume()

	cidr := os.Getenv("PodCIDR")
	if cidr != "" {
		res := getCIDRS(cidr)
		delRoute(res)
	}

	delCNIIface()

	zapDevice()
}

func UnmountPodVolume() {
	fmt.Println("Start umount")
	cmd1 := exec.Command("mount")
	cmd2 := exec.Command("grep", "/var/lib/kubelet")
	var outbuf1, outbuf2 bytes.Buffer
	cmd1.Stdout = &outbuf1
	cmd1.Start()
	cmd1.Wait()
	cmd2.Stdin = &outbuf1
	cmd2.Stdout = &outbuf2
	cmd2.Start()
	cmd2.Wait()
	outputs := strings.Split(outbuf2.String(), "\n")
	for _, l := range outputs {
		if !strings.Contains(l, "kubelet") {
			continue
		}
		path := strings.Fields(l)[2]
		out, err := exec.Command("umount", "-f", path).Output()
		fmt.Println("umount path:", path)
		if err != nil {
			fmt.Println(out, err)
		}
	}
}

func delRoute(nets []string) {
	fmt.Println("Start del router")
	for _, n := range nets {
		fmt.Println("route del net:", n)
		out, err := exec.Command("ip", "r", "flush", n).Output()
		if err != nil {
			fmt.Println(out, err)
		}
	}
}

func delCNIIface() {
	nics := []string{"flannel.1", "cni0"}
	for _, n := range nics {
		out, err := exec.Command("ip", "link", "delete", n).Output()
		if err != nil {
			fmt.Println(out, err)
		}
	}
	out, err := exec.Command("modprobe", "-r", "ipip").Output()
	if err != nil {
		fmt.Println(out, err)
	}
}

func getCIDRS(cidr string) []string {
	res := make([]string, 0)
	cmd1 := exec.Command("ip", "r")
	cmd2 := exec.Command("grep", "/")
	var outbuf1, outbuf2 bytes.Buffer
	cmd1.Stdout = &outbuf1
	cmd1.Start()
	cmd1.Wait()
	cmd2.Stdin = &outbuf1
	cmd2.Stdout = &outbuf2
	cmd2.Start()
	cmd2.Wait()
	outputs := strings.Split(outbuf2.String(), "\n")
	for _, l := range outputs {
		if !strings.Contains(l, "/") {
			continue
		}
		path := strings.Fields(l)
		for _, v := range path {
			if strings.Contains(v, "/") && isContains(cidr, v) {
				res = append(res, v)
			}
		}
	}
	return res
}

func isContains(r, c string) bool {
	_, all, _ := net.ParseCIDR(r)
	it, _, _ := net.ParseCIDR(c)
	return all.Contains(it)
}

func zapDevice() {
	out, err := exec.Command("/bin/sh", "/zap_device.sh").Output()
	if err != nil {
		fmt.Println(out, err)
	}
}
