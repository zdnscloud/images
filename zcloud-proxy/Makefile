
REGISTRY_NAME = zdnscloud
IMAGE_NAME = zcloud-proxy
VERSION = v1.0.3

.PHONY: all container

all: container

container: 
	docker build -t $(REGISTRY_NAME)/$(IMAGE_NAME):${VERSION} ./ --no-cache