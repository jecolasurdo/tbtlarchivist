help: ## this help message
	$(info Available targets)
	@awk '/^[a-zA-Z\-\_0-9]+:/ {                    \
	  nb = sub( /^## /, "", helpMsg );              \
	  if(nb == 0) {                                 \
		helpMsg = $$0;                              \
		nb = sub( /^[^:]*:.* ## /, "", helpMsg );   \
	  }                                             \
	  if (nb)                                       \
		print  $$1 "\t" helpMsg;                    \
	}                                               \
	{ helpMsg = $$0 }'                              \
	$(MAKEFILE_LIST) | column -ts $$'\t' |          \
	grep --color '^[^ ]*'
.PHONY: help

start-all: generate_protobuf start-message-bus start-maria-db ## start all dockers
.PHONY: start-all
	
kill-all: kill-message-bus kill-maria-db ## shut down all dockers and clean up
.PHONY: kill-all

start-message-bus: ## start the rabbitmq docker container
	docker run -d --hostname archivist-mq --name archivist-mq -p 5672:5672  rabbitmq:3.8
.PHONY: start-message-bus

kill-message-bus: ## shut down the rabbitmq docker container and reclaim resources
	docker container kill archivist-mq && \
	docker container prune -f
.PHONY: kill-message-bus

list-queues: ## list queues currently registered in the message bus
	docker exec archivist-mq rabbitmqctl list_queues
.PHONY: list-queues

purge-queues: ## purge the message bus queues
	  docker exec archivist-mq rabbitmqctl purge_queue pending_research
.PHONY: purge-queues

start-maria-db:
	docker run -d --name mariadb -e MYSQL_ALLOW_EMPTY_PASSWORD=true -p 3306:3306 mariadb:10.5
.PHONY: start-maria-db

kill-maria-db: ## shut down the mariadb docker container and reclaim resources
	docker container kill mariadb && \
	docker container prune -f
.PHONY: kill-maria-db

restart: kill-all start-all
.PHONY: restart

bootstrap-maria-db: ## apply db migrations
	echo "create database tbtlarchivist" | mariadb -h 127.0.0.1 -P 3306 -u root	
	flyway migrate
.PHONY: bootstrap-maria-db

generate-protobuf: ## generate concrete implementations of protocol buffer messages.
	go install google.golang.org/protobuf/cmd/protoc-gen-go
	protoc --go_out=. protobuf/contracts.proto
.PHONY: generate-protobuf

generate-mocks: ## generate mocks for testing
	rm -drf ./mocks

	mkdir -p ./mocks/accessors/mock_messagebus
	mockgen -source=pkg/accessors/messagebus/messagebus.go > mocks/accessors/mock_messagebus/mock_messagebus.go

	mkdir -p ./mocks/accessors/mock_messagebus/mock_acknowledger
	mockgen -source=pkg/accessors/messagebus/acknowledger/acknack.go > mocks/accessors/mock_messagebus/mock_acknowledger/mock_acknack.go
.PHONY: generate-mocks

test: generate-mocks ## run unit tests
	go run ./... -race -count=1
.PHONY: test
