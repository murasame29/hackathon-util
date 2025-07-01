build:
	docker compose build

run: build
	docker compose up

rund: build
	docker compose up -d