

start-message-bus:
	docker run -d --hostname archivist-mq --name archivist-mq rabbitmq:3.8
.PHONY: start-message-bus

tear-down:
	docker container kill archivist-mq && \
	docker container prune -f
.PHONY: tear-down