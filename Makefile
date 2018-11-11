CORE_IMAGE_NAME ?= caloriosa/core-dev
CORE_CONTAINER_NAME ?= caloriosa-core
VOLUME ?= $(shell pwd)/config:/config
ENTRYPOINT ?= bin/caloriosa-server

ifeq ($(ENTRYPOINT),bin/caloriosa-server)
	COMMAND ?= -logtostderr -config /config/config.yaml
else
	COMMAND ?= 
endif

build:
	docker build -t $(CORE_IMAGE_NAME) .
run:
	docker run --rm --name $(CORE_CONTAINER_NAME) -ti -p 8080:8080 -v $(VOLUME) -u $(shell id -u):$(shell id -g) --entrypoint "$(ENTRYPOINT)" $(CORE_IMAGE_NAME) $(COMMAND)
dev:
	make run VOLUME=$(shell pwd):/go/src/core ENTRYPOINT=/bin/sh
test:
	make run ENTRYPOINT=/bin/sh COMMAND=./test.sh
shell:
	docker ps | grep $(CORE_CONTAINER_NAME) || (echo "Core container '$(CORE_CONTAINER_NAME)' not running"; exit 1)
	docker exec -ti $(CORE_CONTAINER_NAME) /bin/sh
copy-bin:
	make run VOLUME=$(shell pwd):/mnt ENTRYPOINT="/bin/sh" COMMAND="-c 'cp -r ./bin /mnt'"
