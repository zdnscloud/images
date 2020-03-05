REGISTRY_NAME = zdnscloud
IMAGE_NAME = iscsi-init
BUILD=`date +%FT%T%z`
BRANCH=`git branch | sed -n '/\* /s///p'`
VERSION=latest

.PHONY: all

all: image

image: 
	docker build -t $(REGISTRY_NAME)/$(IMAGE_NAME):${VERSION} . --no-cache
	docker image prune -f

docker:image
	docker push $(REGISTRY_NAME)/$(IMAGE_NAME):${VERSION}

clean-image:
	docker rmi $(REGISTRY_NAME)/$(IMAGE_NAME):${BRANCH}
