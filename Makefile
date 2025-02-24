OWNER ?= ociscloud
PROJECT ?= VirtualRegistryManagement
ABBR ?= vrm
IMAGE_NAME ?= virtual-registry-management
GOVERSION ?= 1.22.4
OS ?= ubuntu
ARCH ?= amd64
PREVERSION ?= 0.0.6
VERSION ?= $(shell cat VERSION)
PWD := $(shell pwd)
GO_PROXY ?= "https://proxy.golang.org,http://proxy.pegasus-cloud.com:8078"

# Release Mode could be dev or prod,
# dev: default, will add commit id to version
# prod: will use version only
RELEASE_MODE ?= dev
COMMIT_ID ?= $(shell git rev-parse --short HEAD)

sed = sed
ifeq ("$(shell uname -s)", "Darwin")	# BSD sed, like MacOS
	sed += -i ''
else	# GNU sed, like LinuxOS
	sed += -i''
endif

ifeq ($(RELEASE_MODE),prod)
    RELEASE_VERSION := $(VERSION)
else
    RELEASE_VERSION := $(VERSION)-$(COMMIT_ID)
endif

.PHONY: go-build
go-build:
	@echo "Build Binary"
	@go build -ldflags="-s -w" -o tmp/$(PROJECT)_$(VERSION)

.PHONY: build
build: go-buildbuild
ifeq ($(OS), ubuntu)
	@sh build/build-debian.sh
else
	@sh build/build-rpm.sh
endif

.PHONY: set-version
set-version:
	@echo "Set Version: $(RELEASE_VERSION)"
	@$(sed) -e'/$(PREVERSION)/{s//$(RELEASE_VERSION)/;:b' -e'n;bb' -e\} $(PWD)/build/$(PROJECT).spec
	@$(sed) -e'/$(PREVERSION)/{s//$(RELEASE_VERSION)/;:b' -e'n;bb' -e\} $(PWD)/constants/common.go
	@$(sed) -e'/$(PREVERSION)/{s//$(RELEASE_VERSION)/;:b' -e'n;bb' -e\} $(PWD)/etc/virtual-registry-management.yaml
	@$(sed) -e'/$(PREVERSION)/{s//$(RELEASE_VERSION)/;:b' -e'n;bb' -e\} $(PWD)/etc/vrm-scheduler.yaml
	@$(sed) -e'/$(PREVERSION)/{s//$(RELEASE_VERSION)/;:b' -e'n;bb' -e\} $(PWD)/Makefile

.PHONY: release
release: 
	@make set-version
	@mkdir -p tmp
	@rm -rf tmp/$(OS)
	@docker run --name build-env -e GOPROXY=$(GO_PROXY) -e GOSUMDB="off" --network=host -v ${PWD}/..:/home -w /home/VirtualRegistryManagement $(OWNER)/golang:$(GOVERSION)-$(OS)-$(ARCH) make OS=$(OS) build
	@docker rm -f build-env
	@mkdir tmp/$(OS)
	@mv tmp/$(PROJECT)* tmp/$(OS)

.PHONY: release-image
release-image: 
	@make set-version
	@rm -rf build/scratch_image/tmp
	@rm -rf tmp/container
	@docker run --name build-env -e GOPROXY=$(GO_PROXY) -e GOSUMDB="off" --network=host -v ${PWD}/..:/home -w /home/VirtualRegistryManagement $(OWNER)/golang:$(GOVERSION)-$(OS)-$(ARCH) make build-container
	@docker rm -f build-env
	@mkdir -p tmp/container
	@docker rmi -f $(OWNER)/$(IMAGE_NAME):$(RELEASE_VERSION)
	@docker build -t $(OWNER)/$(IMAGE_NAME):$(RELEASE_VERSION) build/scratch_image/
	@docker save $(OWNER)/$(IMAGE_NAME):$(RELEASE_VERSION) > tmp/container/$(ABBR)_$(RELEASE_VERSION).image.tar

.PHONY: push-image
push-image:
	@echo "Check Image $(OWNER)/$(IMAGE_NAME):$(RELEASE_VERSION)"
	@docker image inspect $(OWNER)/$(IMAGE_NAME):$(RELEASE_VERSION) --format="image existed"
	@echo "Push Image"
	@docker logout
	@docker login -u $(OWNER) --password-stdin <<< "<DOCKER HUB KEY>"
	@docker image push $(OWNER)/$(IMAGE_NAME):$(RELEASE_VERSION)
	@docker logout

.PHONY: start
start:
	@go env -w GOPROXY=$(GO_PROXY)
	@go env -w GOSUMDB="off"
	@go run main.go -c etc/virtual-registry-management.yaml serve

.PHONY: start-scheduler
start-scheduler:
	@go env -w GOPROXY=$(GO_PROXY)
	@go env -w GOSUMDB="off"
	@go run main.go -c etc/virtual-registry-management.yaml -s etc/vrm-scheduler.yaml scheduler start

.PHONY: init
init:
	@go run main.go -c etc/virtual-registry-management.yaml database sync

.PHONY: start-dev-env
start-dev-env:
	@make start-dev-persistent
	@make start-dev-system
	@make start-dev-service
	
.PHONY: start-dev-service
start-dev-service: docker-compose/service/docker-compose.*.yaml
	@for f in $^; do COMPOSE_IGNORE_ORPHANS=True docker-compose -f $${f} -p "pegasus-service" up -d --no-recreate; done

.PHONY: start-dev-system
start-dev-system: docker-compose/system/docker-compose.*.yaml
	@for f in $^; do COMPOSE_IGNORE_ORPHANS=True docker-compose -f $${f} -p "pegasus-system" up -d --no-recreate; done

.PHONY: start-dev-persistent
start-dev-persistent: docker-compose/persistent/docker-compose.*.yaml
	@for f in $^; do COMPOSE_IGNORE_ORPHANS=True docker-compose -f $${f} -p "pegasus-system" up -d --no-recreate; done

.PHONY: stop-dev-env # Stop and Remove current service only
stop-dev-env:
	COMPOSE_IGNORE_ORPHANS=True docker-compose -f docker-compose/service/docker-compose.${ABBR}.yaml -p "pegasus-service" down
	
.PHONY: stop-dev-all # Stop and Remove all dependency
stop-dev-all:
	@make stop-dev-service
	@make stop-dev-system

.PHONY: purge-dev-all # Stop and Remove all dependency include persistent network and volume
purge-dev-all:
	@make stop-dev-all
	@make clean-dev-persistent

.PHONY: stop-dev-service
stop-dev-service: docker-compose/service/docker-compose.*.yaml
	@for f in $^; do COMPOSE_IGNORE_ORPHANS=True docker-compose -f $${f} -p "pegasus-service" down -v; done

.PHONY: stop-dev-system
stop-dev-system: docker-compose/system/docker-compose.*.yaml
	@for f in $^; do COMPOSE_IGNORE_ORPHANS=True docker-compose -f $${f} -p "pegasus-system" down -v; done

.PHONY: clean-dev-persistent
clean-dev-persistent: docker-compose/persistent/docker-compose.*.yaml
	@for f in $^; do COMPOSE_IGNORE_ORPHANS=True docker-compose -f $${f} -p "pegasus-system" down -v; done

.PHONY: build-container
build-container:
	@go build -o build/scratch_image/tmp/$(PROJECT)
	@sh build/scratch_image/build_scratch_img_env.sh

.PHONY: release-alpine-image
release-alpine-image:
	@make set-version
	@rm -rf build/alpine_image/tmp
	@docker run --name build-env -e GOPROXY=$(GO_PROXY) -e GOSUMDB="off" --network=host -v $(PWD):/home/VirtualRegistryManagement -w /home/VirtualRegistryManagement $(OWNER)/golang:$(GOVERSION)-$(OS)-amd64 make build-alpine-image
	@docker rm -f build-env
	@docker rmi -f $(OWNER)/$(IMAGE_NAME):$(RELEASE_VERSION)
	@docker build -t $(OWNER)/$(IMAGE_NAME):$(RELEASE_VERSION) build/alpine_image/
	@docker save -o tmp/container/$(ABBR)_$(RELEASE_VERSION).image.tar $(OWNER)/$(IMAGE_NAME):$(RELEASE_VERSION)

.PHONY: build-alpine-image
build-alpine-image:
	@go build -o build/alpine_image/tmp/$(PROJECT)
	@sh build/alpine_image/build_alpine_img_env.sh