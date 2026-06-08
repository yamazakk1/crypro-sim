.PHONY: up up-build down clean restart logs seed

up:
	docker-compose up -d

up-build:
	docker-compose up -d --build 

down:
	docker-compose down

clean: 
	docker-compose down -v

restart:
	docker-compose down
	docker-compose up -d

logs:
	docker-compose logs -f

seed:
	go run ./cmd/seed/main.go
	