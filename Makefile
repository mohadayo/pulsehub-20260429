.PHONY: test test-python test-go test-ts lint up down build clean

test: test-python test-go test-ts

test-python:
	cd gateway && pip install -q -r requirements.txt && pytest -v

test-go:
	cd analytics && go test -v ./...

test-ts:
	cd notifier && npm install --silent && npm test

lint: lint-python lint-go lint-ts

lint-python:
	cd gateway && flake8 --max-line-length=120 main.py test_main.py

lint-go:
	cd analytics && go vet ./...

lint-ts:
	cd notifier && npx eslint src/

up:
	docker compose up -d --build

down:
	docker compose down

build:
	docker compose build

clean:
	docker compose down -v --rmi local
