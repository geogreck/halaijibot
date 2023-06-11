all: run


.PHONY: build
build:
	DOCKER_BUILDKIT=1 docker build -t backend .

.PHONY: run
run:
	docker compose -f docker-compose.yml up -d --build backend

.PHONY: lint
lint:
	DOCKER_BUILDKIT=1 docker build . --target lint
