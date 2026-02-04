SHELL := /bin/sh

.PHONY: up up-fresh down build

up:
	docker compose up --build

up-fresh:
	docker compose build --no-cache
	docker compose up

build:
	docker compose build

down:
	docker compose down
