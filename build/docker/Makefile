DOCKER_PS := $(shell docker ps -a -q)
DOCKER_IMAGES := $(shell docker images -a -q)
DOCKER_LOGS := $(shell docker ps|grep coeus|cut -d " " -f1)

.DEFAULT_GOAL:=help

help:
	@echo "Usage: make <target>"
	@echo ""
	@echo "build	build docker image"
	@echo "run	run docker container"
	@echo "logs	show the logs of the container"
	@echo "clean	stop/remove all containers and images"

build:
	@echo "Building docker image for coeus"
	@docker build --progress=plain -t coeus  -f coeus.dockerfile .	

run:
ifndef DOCKER_PS
	@echo "Starting coeus container"
	@docker run --name coeus -h coeus-docker -d  -p 5432:5432 -p 8080:8080  -t coeus
else
	@echo  "Container coeus already running"
endif

logs:
ifdef DOCKER_LOGS
	@docker logs $(DOCKER_LOGS)
endif

clean:
ifdef DOCKER_PS
	@docker rm -f $(DOCKER_PS) 
endif
ifdef DOCKER_IMAGES
	@docker rmi -f $(DOCKER_IMAGES)
endif
	@echo -n ""
