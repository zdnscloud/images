REGISTRY_NAME = zdnscloud
IMAGE_Name = flannel-directrouting
IMAGE_VERSION = v0.0.1

.PHONY: all container

all: container

container: 
	docker build -t $(REGISTRY_NAME)/$(IMAGE_Name):$(IMAGE_VERSION) ./ --no-cache
build:
	CGO_ENABLED=0 GOOS=linux go build -o flanneld
clean:
	rm -f flanneld
	docker rmi $(REGISTRY_NAME)/$(IMAGE_Name):$(IMAGE_VERSION)
