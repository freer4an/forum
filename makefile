# Define the Docker image name and container name
IMAGE_NAME := forum
CONTAINER_NAME := FORUM

.PHONY: all build run stop clean

all: build run

build:
	docker build -t $(IMAGE_NAME) .

run:
	docker run -dp 8000:8000 --name $(CONTAINER_NAME) $(IMAGE_NAME)

stop:
	docker stop $(CONTAINER_NAME)

clean:
	docker rm $(CONTAINER_NAME)
	docker rmi $(IMAGE_NAME)
