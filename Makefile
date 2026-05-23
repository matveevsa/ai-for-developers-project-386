.PHONY: build dev down spec check clean

build:
	docker compose build

dev:
	docker compose up -d

down:
	docker compose down

spec:
	docker compose --profile compile run --rm spec

check: spec
	docker compose run --rm frontend npm run build

clean:
	docker compose down -v
