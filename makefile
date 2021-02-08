start-all: start-message-bus start-maria-db
.PHONY: start-all
	
kill-all: kill-message-bus kill-maria-db
.PHONY: kill-all

start-message-bus:
	docker run -d --hostname archivist-mq --name archivist-mq -p 5672:5672  rabbitmq:3.8
.PHONY: start-message-bus

kill-message-bus:
	docker container kill archivist-mq && \
	docker container prune -f
.PHONY: kill-message-bus

start-maria-db:
	docker run -d --name mariadb -e MYSQL_ALLOW_EMPTY_PASSWORD=true -p 3306:3306 mariadb:10.5
.PHONY: start-maria-db

kill-maria-db:
	docker container kill mariadb && \
	docker container prune -f
.PHONY: kill-maria-db

restart: kill-all start-all
.PHONY: restart