PRODUCER_BUILD_FILE=producer
CONSUMER_BUILD_FILE=consumer
LOGSEARCH_BUILD_FILE=logsearch
PRODUCER_DOCKER_FILE=producer-image
CONSUMER_DOCKER_FILE=consumer-image
LOGSEARCH_DOCKER_FILE=logsearch-image

all: build-producer build-consumer build-logsearch

docker-compose-up:
	docker-compose up -d

build-producer:
	go build -o ./bin/$(PRODUCER_BUILD_FILE) ./producer
	chmod +x ./bin/$(PRODUCER_BUILD_FILE)

build-consumer:
	go build -o ./bin/$(CONSUMER_BUILD_FILE) ./consumer
	chmod +x ./bin/$(CONSUMER_BUILD_FILE)

build-logsearch:
	go build -o ./bin/$(LOGSEARCH_BUILD_FILE) ./logsearch
	chmod +x ./bin/$(LOGSEARCH_BUILD_FILE)

run-producer:
	./bin/$(PRODUCER_BUILD_FILE)

run-consumer:
	./bin/$(CONSUMER_BUILD_FILE)

run-logsearch:
	./bin/$(LOGSEARCH_BUILD_FILE)

docker-build-producer:
	docker build -t $(PRODUCER_DOCKER_FILE) -f Dockerfile.producer .

docker-build-consumer:
	docker build -t $(CONSUMER_DOCKER_FILE) -f Dockerfile.consumer .

docker-build-logsearch:
	docker build -t $(LOGSEARCH_DOCKER_FILE) -f Dockerfile.logsearch .

docker-run-producer:
	docker run --rm -p 3000:3000 $(PRODUCER_DOCKER_FILE)

docker-run-consumer:
	docker run --rm $(CONSUMER_DOCKER_FILE)

docker-run-logsearch:
	docker run --rm $(LOGSEARCH_DOCKER_FILE)

clean:
	rm -f ./bin/$(PRODUCER_BUILD_FILE) ./bin/$(CONSUMER_BUILD_FILE) ./bin/$(LOGSEARCH_BUILD_FILE)
	docker-compose down

docker-delete:
	docker rmi -f $(PRODUCER_DOCKER_FILE):latest
	docker rmi -f producer:latest
	docker rmi -f $(CONSUMER_DOCKER_FILE):latest
	docker rmi -f consumer:latest
	docker rmi -f $(LOGSEARCH_DOCKER_FILE):latest
	docker rmi -f logsearch:latest