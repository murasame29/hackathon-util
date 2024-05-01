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