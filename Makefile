TARGET = certserv
VERSION = $(shell git describe --tags)
BUILD = $(shell date +"%F_%T_%Z")
LEVEL = $(shell git log --pretty=format:"%H" --name-status HEAD^..HEAD | head -1)
DOCKERIMAGE = $(TARGET):$(VERSION)
DOCKERFILE = Dockerfile

all: build

build:
	go build -o $(TARGET)

test:
	go test

image:
	docker build -f $(DOCKERFILE) -t $(DOCKERIMAGE) .

clean:
	go clean
	rm -f $(TARGET) *~

.PHONY: test clean 
