package main

type component struct {
	Name   string
	Images map[string]string
}

var (
	images = map[string][]component{
		"v2.0": []component{
			component{
				Name:   "zke",
				Images: zkeImagesV20,
			},
			component{
				Name:   "zcloud",
				Images: zcloudImagesV20,
			},
			component{
				Name:   "monitor",
				Images: monitorImagesV20,
			},
			component{
				Name:   "registry",
				Images: registryImagesV20,
			},
			component{
				Name:   "storage",
				Images: storageImagesV20,
			},
		}}

	zkeImagesV20 = map[string]string{
		"Etcd":                      "zdnscloud/coreos-etcd:v3.3.10",
		"Kubernetes":                "zdnscloud/hyperkube:v1.13.10",
		"Alpine":                    "zdnscloud/zke-tools:v0.1.40",
		"NginxProxy":                "zdnscloud/zke-tools:v0.1.40",
		"CertDownloader":            "zdnscloud/zke-tools:v0.1.40",
		"KubernetesServicesSidecar": "zdnscloud/zke-tools:v0.1.40",
		"Flannel":                   "zdnscloud/coreos-flannel:v0.10.0",
		"FlannelCNI":                "zdnscloud/coreos-flannel-cni:v0.3.0",
		"CalicoNode":                "zdnscloud/calico-node:v3.4.0",
		"CalicoCNI":                 "zdnscloud/calico-cni:v3.4.0",
		"CalicoCtl":                 "zdnscloud/calico-ctl:v2.0.0",
		"PodInfraContainer":         "zdnscloud/pause-amd64:3.1",
		"Ingress":                   "zdnscloud/nginx-ingress-controller:0.21.0",
		"IngressBackend":            "zdnscloud/nginx-ingress-controller-defaultbackend:1.4",
		"CoreDNS":                   "zdnscloud/coredns:1.2.6",
		"CoreDNSAutoscaler":         "zdnscloud/cluster-proportional-autoscaler-amd64:1.0.0",
		"ClusterAgent":              "zdnscloud/cluster-agent:v2.5",
		"NodeAgent":                 "zdnscloud/node-agent:v1.2",
		"MetricsServer":             "zdnscloud/metrics-server-amd64:v0.3.3",
		"ZKERemover":                "zdnscloud/zke-remove:v0.5",
		"StorageOperator":           "zdnscloud/storage-operator:v1.8",
		"ZcloudShell":               "zdnscloud/kubectl:v1.13.1",
	}

	zcloudImagesV20 = map[string]string{
		"ZcloudProxy": "zdnscloud/zcloud-proxy:v1.0.2",
		"Zcloud":      "zdnscloud/zcloud:v2.0",
	}

	monitorImagesV20 = map[string]string{
		"GrafanaSideCar":           "kiwigrid/k8s-sidecar:0.0.18",
		"KubeStateMetrics":         "quay.io/coreos/kube-state-metrics:v1.7.2",
		"NodeExporter":             "quay.io/prometheus/node-exporter:v0.18.0",
		"Grafana":                  "zdnscloud/grafana:6.2.5",
		"AlertManager":             "zdnscloud/prometheus-alertmanager:v0.17.0",
		"PrometheusConfigReloader": "zdnscloud/prometheus-config-reloader:v0.31.1",
		"PrometheusOperator":       "zdnscloud/prometheus-operator:v0.31.1",
		"Prometheus":               "zdnscloud/prometheus:v2.10.0",
		"GrafanaInit":              "busybox:1.30.0",
	}

	registryImagesV20 = map[string]string{
		"Chartmuseum": "zdnscloud/goharbor-chartmuseum-photon:v0.8.1-v1.8.1",
		"Clair":       "zdnscloud/goharbor-clair-photon:v2.0.8-v1.8.1",
		"Core":        "zdnscloud/goharbor-harbor-core:v1.8.1",
		"DB":          "zdnscloud/goharbor-harbor-db:v1.8.1",
		"Jobservice":  "zdnscloud/goharbor-harbor-jobservice:v1.8.1",
		"Portal":      "zdnscloud/goharbor-harbor-portal:v1.8.1",
		"Registryctl": "zdnscloud/goharbor-harbor-registryctl:v1.8.1",
		"Redis":       "zdnscloud/goharbor-redis-photon:v1.8.1",
		"Registry":    "zdnscloud/goharbor-registry-photon:v2.7.1-patch-2819-v1.8.1",
	}

	storageImagesV20 = map[string]string{
		"LvmCsi":                 "zdnscloud/lvmcsi:v0.6",
		"Lvmd":                   "zdnscloud/lvmd:v0.5",
		"CsiAttacher":            "quay.io/k8scsi/csi-attacher:v1.0.1",
		"CsiNodeDriverRegistrar": "quay.io/k8scsi/csi-node-driver-registrar:v1.0.2",
		"CsiProvisioner":         "quay.io/k8scsi/csi-provisioner:v1.0.1",
		"Busybox":                "busybox:1.31.0",
		"CephDaemon":             "ceph/daemon:latest-mimic",
		"CephDaemonLatest":       "ceph/daemon:latest",
		"CephfsPlugin":           "quay.io/cephcsi/cephfsplugin:v1.0.0",
	}
)