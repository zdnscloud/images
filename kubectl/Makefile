
REGISTRY_NAME = zdnscloud
IMAGE_NAME = kubectl
VERSION = v1.13

.PHONY: all container

all: container

container: 
	docker build -t $(REGISTRY_NAME)/$(IMAGE_NAME):${VERSION} ./ --no-cache