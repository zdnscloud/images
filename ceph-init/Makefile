REGISTRY_NAME = zdnscloud
IMAGE_Name = ceph-init
IMAGE_VERSION = v0.3

.PHONY: all container

all: container

container: 
	docker build -t $(REGISTRY_NAME)/$(IMAGE_Name):$(IMAGE_VERSION) ./ --no-cache
