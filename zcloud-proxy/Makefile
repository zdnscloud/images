REGISTRY_NAME = zdnscloud
IMAGE_NAME = zcloud-proxy
VERSION = v1.0.4


.PHONY: docker container

docker: container
	docker push $(REGISTRY_NAME)/$(IMAGE_NAME):${VERSION}
container: 
	docker build -t $(REGISTRY_NAME)/$(IMAGE_NAME):${VERSION} ./ --no-cache --build-arg goproxy=${GOPROXY}
	