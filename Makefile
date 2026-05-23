.PHONY: build dev down spec check clean test

build:
	docker compose build

dev:
	docker compose up -d

down:
	docker compose down

spec:
	docker compose --profile compile run --rm spec

check: spec
	docker compose build app

test:
	go test ./... -count=1

test-verbose:
	go test ./... -v -count=1

clean:
	docker compose down -v
