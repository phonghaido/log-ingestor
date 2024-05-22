PRODUCER_BUILD_FILE=producer
CONSUMER_BUILD_FILE=consumer
PRODUCER_DOCKER_FILE=producer-image
CONSUMER_DOCKER_FILE=consumer-image

all: build-producer build-consumer

docker-compose-up:
	docker-compose up -d

build-producer:
	go build -o ./bin/$(PRODUCER_BUILD_FILE) ./producer
	chmod +x ./bin/$(PRODUCER_BUILD_FILE)

build-consumer:
	go build -o ./bin/$(CONSUMER_BUILD_FILE) ./consumer
	chmod +x ./bin/$(CONSUMER_BUILD_FILE)

run-producer:
	./bin/$(PRODUCER_BUILD_FILE)

run-consumer:
	./bin/$(CONSUMER_BUILD_FILE)

docker-build-producer:
	docker build -t $(PRODUCER_DOCKER_FILE) -f Dockerfile.producer .

docker-build-consumer:
	docker build -t $(CONSUMER_DOCKER_FILE) -f Dockerfile.consumer .

docker-run-producer:
	docker run --rm -p 3000:3000 $(PRODUCER_DOCKER_FILE)

docker-run-consumer:
	docker run --rm $(CONSUMER_DOCKER_FILE)

clean:
	rm -f ./bin/$(PRODUCER_BUILD_FILE) ./bin/$(CONSUMER_BUILD_FILE)
	docker-compose down