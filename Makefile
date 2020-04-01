APP_NAME=patches
REGISTRY?=343660461351.dkr.ecr.us-east-2.amazonaws.com
TAG?=latest
PATCHES_HEIMDALL_SERVER?=localhost
PATCHES_ETHER_SERVER?=localhost:8082
PATCHES_KAFKA_SERVER?=localhost:9092
PATCHES_KAFKA_TOPIC?=updates
PATCHES_DB_HOST?=localhost
PATCHES_DB_PORT?=5432
PATCHES_DB_USERNAME?=postgres
PATCHES_DB_PASSWORD?=patches
HELP_FUNC = \
    %help; \
    while(<>) { \
        if(/^([a-z0-9_-]+):.*\#\#(?:@(\w+))?\s(.*)$$/) { \
            push(@{$$help{$$2 // 'targets'}}, [$$1, $$3]); \
        } \
    }; \
    print "usage: make [target]\n\n"; \
    for ( sort keys %help ) { \
        print "$$_:\n"; \
        printf("  %-20s %s\n", $$_->[0], $$_->[1]) for @{$$help{$$_}}; \
        print "\n"; \
    }

.PHONY: help
help: 				## show options and their descriptions
	@perl -e '$(HELP_FUNC)' $(MAKEFILE_LIST)

all:  				## clean the working environment, build and test the packages, and then build the docker image
all: clean test docker

tmp: 				## create tmp/
	if [ -d "./tmp" ]; then rm -rf ./tmp; fi
	mkdir tmp

build: tmp 			## build the app binaries
	go build -o ./tmp ./...

test: build 		## build and test the module packages
	go test ./...

run: build 			## build and run the app binaries
	export PATCHES_HEIMDALL_SERVER="${PATCHES_HEIMDALL_SERVER}" && \
		export PATCHES_ETHER_SERVER="${PATCHES_ETHER_SERVER}" && \
		export PATCHES_KAFKA_SERVER="${PATCHES_KAFKA_SERVER}" && \
		export PATCHES_KAFKA_TOPIC="${PATCHES_KAFKA_TOPIC}" && \
		export PATCHES_DB_HOST="${PATCHES_DB_HOST}" && \
		export PATCHES_DB_PORT="${PATCHES_DB_PORT}" && \
		export PATCHES_DB_USERNAME="${PATCHES_DB_USERNAME}" && \
		export PATCHES_DB_PASSWORD="${PATCHES_DB_PASSWORD}" && \
		./tmp/app

docker: tmp 		## build the docker image
	docker build -t $(REGISTRY)/$(APP_NAME):$(TAG) .

docker-run: docker 	## start the built docker image in a container
	docker run -it --rm -p 80:80 --name $(APP_NAME) $(REGISTRY)/$(APP_NAME):$(TAG)

docker-push: tmp docker
	docker push $(REGISTRY)/$(APP_NAME):$(TAG)

.PHONY: clean
clean: 				## remove tmp/ and old docker images
	rm -rf tmp
ifneq ("$(shell docker container list -a | grep $(APP_NAME))", "")
	docker rm -f $(APP_NAME)
endif
	docker system prune
ifneq ("$(shell docker images | grep $(APP_NAME) | awk '{ print $$3; }')", "") 
	docker images | grep $(APP_NAME) | awk '{ print $$3; }' | xargs -I {} docker rmi -f {}
endif
