create_role:
	go run cmd/create_role/main.go -e .env

delete_role:
	go run cmd/delete_role/main.go -e .env

create_channel:
	go run cmd/create_channel/main.go -e .env

delete_channel:
	go run cmd/delete_channel/main.go -e .env

bind_role:
	go run cmd/bind_role/main.go -e .env

sync_role:
	go run cmd/delete_role/main.go -e .env
	go run cmd/create_role/main.go -e .env
	go run cmd/bind_role/main.go -e .env

build:
	docker compose build

run: build
	docker compose up

rund: build
	docker compose up -d

build-sheetless: 
	docker compose -f sheetless.compose.yaml build

run-sheetless: build-sheetless
	docker compose -f sheetless.compose.yaml up --build

rund-sheetless: build-sheetless
	docker compose -f sheetless.compose.yaml up -d