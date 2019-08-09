REGISTRY_NAME = zdnscloud
IMAGE_Name = zke-remove
IMAGE_VERSION = v0.5

.PHONY: all container

all: container

container: 
	docker build -t $(REGISTRY_NAME)/$(IMAGE_Name):$(IMAGE_VERSION) ./ --no-cache
