
REGISTRY_NAME = zdnscloud
IMAGE_NAME = zke-tools
VERSION = v0.1.39

.PHONY: all container

all: container

container: 
	docker build -t $(REGISTRY_NAME)/$(IMAGE_NAME):${VERSION} ./ --no-cache