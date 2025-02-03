.PHONY: run
run:
	docker-compose -f deploy/docker-compose.yml up --build

.PHONY: down
down:
	docker-compose -f deploy/docker-compose.yml down
